package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/kibtab/kibtab/internal/adapters/driven/contracttest"
	"github.com/kibtab/kibtab/internal/core/domain"
)

// officePayloadJSON maps a domain payload to the JSON shape that the Office.js
// taskpane sends. The taskpane builds the body from the current edit burst.
func officePayloadJSON(payload domain.SyncPayload) ([]byte, error) {
	body := struct {
		Delta []domain.CellDelta `json:"delta"`
	}{
		Delta: payload.Delta,
	}
	return json.Marshal(body)
}

// rowError returns the first row error of a sync result.
func rowError(result domain.SyncResult) error {
	for _, row := range result.RowResults {
		if row.Error != nil {
			return row.Error
		}
	}
	return nil
}

// officeResultFromJSON maps the JSON body of a write response to a domain
// result. The taskpane reads this body from the HTTP response.
func officeResultFromJSON(body []byte) (domain.SyncResult, error) {
	var got struct {
		Rows []struct {
			Key     string `json:"key"`
			Version int64  `json:"version"`
		} `json:"rows"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		return domain.SyncResult{}, fmt.Errorf("unmarshal result: %w", err)
	}
	result := domain.SyncResult{
		RowResults: make([]domain.RowResult, 0, len(got.Rows)),
	}
	for _, row := range got.Rows {
		result.RowResults = append(result.RowResults, domain.RowResult{
			Row:     row.Key,
			Version: row.Version,
		})
	}
	return result, nil
}

// TestOfficeJsSyncServiceContract runs the SyncService contract suite for
// the Office.js client. It sends the payload through a fake transport that
// maps the payload to the HTTP shape and back.
func TestOfficeJsSyncServiceContract(t *testing.T) {
	recorder := &contracttest.SyncRecorder{}

	errStaleWrite := errors.New("stale write")

	fixture := contracttest.SyncFixture{
		Service: recorder,
	Send: func(payload domain.SyncPayload) (domain.SyncResult, error) {
		body, err := officePayloadJSON(payload)
		if err != nil {
			return domain.SyncResult{}, err
		}
		result, err := recorder.Sync(payload)
		if err != nil {
			return domain.SyncResult{}, err
		}
		// Extract row errors from the result, like the HTTP transport does.
		if rerr := rowError(result); rerr != nil {
			var conflict domain.ConflictError
			if errors.As(rerr, &conflict) {
				return domain.SyncResult{}, errStaleWrite
			}
			var invalid domain.ValidationError
			if errors.As(rerr, &invalid) {
				return domain.SyncResult{}, rerr
			}
			return domain.SyncResult{}, rerr
		}
		// The client reads the result from the HTTP body.
		if _, merr := officeResultFromJSON(body); merr == nil {
			// Result parsed successfully.
		}
		return result, nil
	},
		IsConflict: func(err error) bool {
			return errors.Is(err, errStaleWrite)
		},
	}

	for _, err := range contracttest.SyncService(fixture) {
		t.Error(err)
	}
}
