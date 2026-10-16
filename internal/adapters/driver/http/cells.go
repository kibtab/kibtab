package http

import (
	nethttp "net/http"

	"github.com/go-chi/chi/v5"
)

// cellsRequest is the JSON body of POST /v1/tables/{table}/cells.
type cellsRequest struct {
	Delta []deltaJSON `json:"delta"`
}

// addCells answers POST /v1/tables/{table}/cells.
func (h handlers) addCells(w nethttp.ResponseWriter, r *nethttp.Request) {
	table := chi.URLParam(r, "table")

	var body cellsRequest
	if err := decodeJSON(r, &body); err != nil {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			"the JSON body is invalid")
		return
	}

	h.writeSync(w, table, body.Delta)
}
