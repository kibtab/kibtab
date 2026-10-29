package http

import (
	"errors"
	"fmt"
	nethttp "net/http"
	"strings"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// deltaJSON is one cell change in the JSON of a cell request.
type deltaJSON struct {
	Row     string `json:"row"`
	Field   string `json:"field"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
}

// syncRowJSON is one row result in the JSON of a write response.
type syncRowJSON struct {
	Key     string `json:"key"`
	Version int64  `json:"version"`
}

// syncJSON is the JSON body of a write response.
type syncJSON struct {
	Rows []syncRowJSON `json:"rows"`
}

// headerUser is the request header that names the user of a write.
const headerUser = "X-Kibtab-User"

// unknownUser is the audit user of a write that names no user.
const unknownUser = "unknown"

// writeSync checks the delta against the table and runs the sync. It
// writes the response for the result. It reads the user of the write
// from the request header.
func (h handlers) writeSync(w nethttp.ResponseWriter, r *nethttp.Request,
	table string, delta []deltaJSON) {
	meta, err := h.tables.Metadata(r.Context(), table)
	if err != nil {
		h.log.Error("read the table metadata", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the table metadata failed")
		return
	}
	field, known := unknownField(delta, meta.Fields)
	if !known {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			fmt.Sprintf("the field %q is not in the table", field))
		return
	}

	payload := toDeltaPayload(table, delta)
	payload.User = requestUser(r)

	result, err := h.syncService.Sync(payload)
	if err != nil {
		h.writeSyncError(w, err)
		return
	}
	if err := rowError(result); err != nil {
		h.writeSyncError(w, err)
		return
	}

	h.writeJSON(w, nethttp.StatusOK, toSyncJSON(result))
}

// writeSyncError writes the response for a failed sync.
func (h handlers) writeSyncError(w nethttp.ResponseWriter, err error) {
	var conflict domain.ConflictError
	if errors.As(err, &conflict) {
		h.writeError(w, nethttp.StatusConflict, codeConflict, err.Error())
		return
	}

	var invalid domain.ValidationError
	if errors.As(err, &invalid) {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			err.Error())
		return
	}

	h.log.Error("run the sync", "error", err)
	h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
		"the row write failed")
}

// requestUser returns the user of one write request. It returns the
// unknown user when the header is empty or when it holds spaces only.
func requestUser(r *nethttp.Request) string {
	user := strings.TrimSpace(r.Header.Get(headerUser))
	if user == "" {
		return unknownUser
	}
	return user
}

// unknownField returns the first delta field outside the table fields.
// The boolean reports that every field is a table field.
func unknownField(delta []deltaJSON, fields []string) (string, bool) {
	known := fieldSet(fields)
	for _, cell := range delta {
		if !known[cell.Field] {
			return cell.Field, false
		}
	}
	return "", true
}

// toDeltaPayload maps the cell changes of one request to a sync payload.
func toDeltaPayload(table string, delta []deltaJSON) domain.SyncPayload {
	cells := make([]domain.CellDelta, 0, len(delta))
	for _, cell := range delta {
		cells = append(cells, domain.CellDelta{
			Row:     cell.Row,
			Field:   cell.Field,
			Value:   cell.Value,
			Version: cell.Version,
		})
	}
	return domain.SyncPayload{Table: table, Delta: cells}
}

// toSyncJSON maps a sync result to the JSON of a write response.
func toSyncJSON(result domain.SyncResult) syncJSON {
	body := syncJSON{Rows: make([]syncRowJSON, 0, len(result.RowResults))}
	for _, row := range result.RowResults {
		body.Rows = append(body.Rows, syncRowJSON{
			Key:     row.Row,
			Version: row.Version,
		})
	}
	return body
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
