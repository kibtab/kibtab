package http

import (
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

// rowJSON is one row in the JSON of a row request or a row response.
type rowJSON struct {
	Key   string     `json:"key"`
	Cells []cellJSON `json:"cells"`
}

// rowList is the JSON body of GET /v1/tables/{table}/rows.
type rowList struct {
	Rows []rowJSON `json:"rows"`
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

	var body rowJSON
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

	if err := h.rows.Write(table, body.Key, toCellValues(body.Cells)); err != nil {
		h.log.Error("write the row", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the row write failed")
		return
	}

	h.writeJSON(w, nethttp.StatusOK, body)
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

// toCellValues maps the JSON cells to domain values.
func toCellValues(cells []cellJSON) []domain.CellValue {
	values := make([]domain.CellValue, 0, len(cells))
	for _, cell := range cells {
		values = append(values, domain.CellValue{
			Field: cell.Field,
			Value: cell.Value,
		})
	}
	return values
}

// fieldSet returns the set of the field names of one table.
func fieldSet(fields []string) map[string]bool {
	set := make(map[string]bool, len(fields))
	for _, field := range fields {
		set[field] = true
	}
	return set
}
