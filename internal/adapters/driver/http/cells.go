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
	if err := decodeJSON(w, r, &body); err != nil {
		h.writeDecodeError(w, err)
		return
	}

	h.writeSync(w, r, table, body.Delta)
}
