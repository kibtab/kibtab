package http

import (
	nethttp "net/http"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// tableResponse is one table in the JSON of the table list.
type tableResponse struct {
	Name    string   `json:"name"`
	Fields  []string `json:"fields"`
	Version int64    `json:"version"`
}

// tableList is the JSON body of GET /v1/tables.
type tableList struct {
	Tables []tableResponse `json:"tables"`
}

// listTables answers GET /v1/tables.
func (h handlers) listTables(w nethttp.ResponseWriter, r *nethttp.Request) {
	tables, err := h.tables.List()
	if err != nil {
		h.log.Error("list the tables", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the table list failed")
		return
	}

	body := tableList{Tables: make([]tableResponse, 0, len(tables))}
	for _, table := range tables {
		body.Tables = append(body.Tables, toTableResponse(table))
	}
	h.writeJSON(w, nethttp.StatusOK, body)
}

// toTableResponse maps one table to its JSON shape.
func toTableResponse(table domain.TableMetadata) tableResponse {
	return tableResponse{
		Name:    table.Name,
		Fields:  table.Fields,
		Version: table.Version,
	}
}
