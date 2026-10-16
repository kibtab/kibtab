package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// cellResponse is one cell in the JSON of a row.
type cellResponse struct {
	Field string `json:"field"`
	Value string `json:"value"`
}

// rowResponse is one row in the JSON of the row list.
type rowResponse struct {
	Key   string         `json:"key"`
	Cells []cellResponse `json:"cells"`
}

// rowList is the JSON body of GET /v1/tables/{table}/rows.
type rowList struct {
	Rows []rowResponse `json:"rows"`
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

// groupRows groups the cells of one table into rows. A cell outside the
// table fields starts a row, and its value is the row key.
func groupRows(cells []domain.CellValue, fields []string) []rowResponse {
	known := fieldSet(fields)
	rows := make([]rowResponse, 0)
	for _, cell := range cells {
		if !known[cell.Field] {
			rows = append(rows, rowResponse{
				Key:   cell.Value,
				Cells: []cellResponse{},
			})
			continue
		}
		if len(rows) == 0 {
			continue
		}
		last := &rows[len(rows)-1]
		last.Cells = append(last.Cells,
			cellResponse{Field: cell.Field, Value: cell.Value})
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
