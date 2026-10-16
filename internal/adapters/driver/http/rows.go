package http

import (
	"errors"
	"fmt"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// cellJSON is one cell in the JSON of a row. A request and a response
// use it.
type cellJSON struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// rowJSON is one row in the JSON of the row list.
type rowJSON struct {
	Key   string     `json:"key"`
	Cells []cellJSON `json:"cells"`
}

// rowList is the JSON body of GET /v1/tables/{table}/rows.
type rowList struct {
	Rows []rowJSON `json:"rows"`
}

// rowWrite is the JSON body of POST /v1/tables/{table}/rows.
type rowWrite struct {
	Version int64      `json:"version"`
	Key     string     `json:"key"`
	Cells   []cellJSON `json:"cells"`
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

// listRows answers GET /v1/tables/{table}/rows.
func (h handlers) listRows(w nethttp.ResponseWriter, r *nethttp.Request) {
	table := chi.URLParam(r, "table")

	meta, err := h.tables.Metadata(table)
	if err != nil {
		h.log.Error("read the table metadata", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the table metadata failed")
		return
	}

	cells, err := h.rows.Read(table)
	if err != nil {
		h.log.Error("read the rows", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the row list failed")
		return
	}

	h.writeJSON(w, nethttp.StatusOK,
		rowList{Rows: groupRows(cells, meta.Fields)})
}

// addRow answers POST /v1/tables/{table}/rows.
func (h handlers) addRow(w nethttp.ResponseWriter, r *nethttp.Request) {
	table := chi.URLParam(r, "table")

	var body rowWrite
	if err := decodeJSON(r, &body); err != nil {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			"the JSON body is invalid")
		return
	}
	if body.Key == "" {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			"the row key is empty")
		return
	}

	meta, err := h.tables.Metadata(table)
	if err != nil {
		h.log.Error("read the table metadata", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the table metadata failed")
		return
	}
	field, known := unknownField(body.Cells, meta.Fields)
	if !known {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			fmt.Sprintf("the field %q is not in the table", field))
		return
	}

	result, err := h.syncService.Sync(toSyncPayload(table, body))
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

// groupRows groups the cells of one table into rows. A cell outside the
// table fields starts a row, and its value is the row key.
func groupRows(cells []domain.CellValue, fields []string) []rowJSON {
	known := fieldSet(fields)
	rows := make([]rowJSON, 0)
	for _, cell := range cells {
		if !known[cell.Field] {
			rows = append(rows, rowJSON{
				Key:   cell.Value,
				Cells: []cellJSON{},
			})
			continue
		}
		if len(rows) == 0 {
			continue
		}
		last := &rows[len(rows)-1]
		last.Cells = append(last.Cells,
			cellJSON{Field: cell.Field, Value: cell.Value})
	}
	return rows
}

// unknownField returns the first cell field outside the table fields.
// The boolean reports that every field is a table field.
func unknownField(cells []cellJSON, fields []string) (string, bool) {
	known := fieldSet(fields)
	for _, cell := range cells {
		if !known[cell.Field] {
			return cell.Field, false
		}
	}
	return "", true
}

// toSyncPayload maps one row write to a sync payload.
func toSyncPayload(table string, body rowWrite) domain.SyncPayload {
	delta := make([]domain.CellDelta, 0, len(body.Cells))
	for _, cell := range body.Cells {
		delta = append(delta, domain.CellDelta{
			Row:     body.Key,
			Field:   cell.Field,
			Value:   cell.Value,
			Version: body.Version,
		})
	}
	return domain.SyncPayload{Table: table, Delta: delta}
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

// fieldSet returns the set of the field names of one table.
func fieldSet(fields []string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, field := range fields {
		set[field] = true
	}
	return set
}
