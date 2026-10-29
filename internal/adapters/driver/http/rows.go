package http

import (
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

// listRows answers GET /v1/tables/{table}/rows.
func (h handlers) listRows(w nethttp.ResponseWriter, r *nethttp.Request) {
	table := chi.URLParam(r, "table")

	meta, err := h.tables.Metadata(r.Context(), table)
	if err != nil {
		h.log.Error("read the table metadata", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the table metadata failed")
		return
	}

	cells, err := h.rows.Read(r.Context(), table)
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
	if err := decodeJSON(w, r, &body); err != nil {
		h.writeDecodeError(w, err)
		return
	}
	if body.Key == "" {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			"the row key is empty")
		return
	}

	h.writeSync(w, r, table, rowDelta(body))
}

// rowDelta expands one row write into one cell change for each cell.
func rowDelta(body rowWrite) []deltaJSON {
	delta := make([]deltaJSON, 0, len(body.Cells))
	for _, cell := range body.Cells {
		delta = append(delta, deltaJSON{
			Row:     body.Key,
			Field:   cell.Field,
			Value:   cell.Value,
			Version: body.Version,
		})
	}
	return delta
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

// fieldSet returns the set of the field names of one table.
func fieldSet(fields []string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, field := range fields {
		set[field] = true
	}
	return set
}
