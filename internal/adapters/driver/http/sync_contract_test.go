package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	nethttp "net/http"
	"net/http/httptest"
	"testing"

	"github.com/kibtab/kibtab/internal/adapters/driven/contracttest"
	"github.com/kibtab/kibtab/internal/core/domain"
)

// errStaleWrite tells the caller that a write used a stale version.
var errStaleWrite = errors.New("stale write")

func TestSyncServiceContract(t *testing.T) {
	recorder := &contracttest.SyncRecorder{}
	handler := NewHandler(usersTable(), fakeRows{}, &fakeAudit{}, recorder, discardLogger())

	fixture := contracttest.SyncFixture{
		Service: recorder,
		Send:    sendPayload(handler),
		IsConflict: func(err error) bool {
			return errors.Is(err, errStaleWrite)
		},
	}
	for _, err := range contracttest.SyncService(fixture) {
		t.Error(err)
	}
}

// sendPayload sends one payload through the cells route and it maps the
// response back to a sync result.
func sendPayload(handler nethttp.Handler) func(domain.SyncPayload) (domain.SyncResult, error) {
	return func(payload domain.SyncPayload) (domain.SyncResult, error) {
		body, err := json.Marshal(cellsRequest{Delta: requestDelta(payload.Delta)})
		if err != nil {
			return domain.SyncResult{}, err
		}
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(nethttp.MethodPost,
			"/v1/tables/"+payload.Table+"/cells", bytes.NewReader(body)))

		switch recorder.Code {
		case nethttp.StatusOK:
			return syncResult(recorder.Body.Bytes())
		case nethttp.StatusConflict:
			return domain.SyncResult{}, staleWriteError(recorder.Body.Bytes())
		default:
			return domain.SyncResult{}, responseError(recorder.Body.Bytes())
		}
	}
}

// requestDelta maps the cell changes of a payload to the request shape.
func requestDelta(delta []domain.CellDelta) []deltaJSON {
	out := make([]deltaJSON, 0, len(delta))
	for _, cell := range delta {
		out = append(out, deltaJSON{
			Row:     cell.Row,
			Field:   cell.Field,
			Value:   cell.Value,
			Version: cell.Version,
		})
	}
	return out
}

// syncResult maps the body of a write response to a sync result.
func syncResult(body []byte) (domain.SyncResult, error) {
	var got syncJSON
	if err := json.Unmarshal(body, &got); err != nil {
		return domain.SyncResult{}, err
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

// staleWriteError maps a conflict response to an error of the caller.
func staleWriteError(body []byte) error {
	var response ErrorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	return fmt.Errorf("%w: %s", errStaleWrite, response.Message)
}

// responseError maps a failure response to an error of the caller.
func responseError(body []byte) error {
	var response ErrorResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}
	return errors.New(response.Message)
}
