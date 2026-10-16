package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kibtab/kibtab/internal/core/ports"
)

// handlers holds the ports and the logger of the HTTP transport.
type handlers struct {
	tables ports.TableRegistry
	log    *slog.Logger
}

// NewHandler builds the HTTP routes of the transport.
func NewHandler(tables ports.TableRegistry, log *slog.Logger) nethttp.Handler {
	h := handlers{tables: tables, log: log}

	router := chi.NewRouter()
	router.Get("/v1/tables", h.listTables)
	return router
}

// writeJSON writes value as a JSON response. It logs a write failure.
func (h handlers) writeJSON(w nethttp.ResponseWriter, status int, value any) {
	if err := encodeJSON(w, status, value); err != nil {
		h.log.Error("write the response", "error", err)
	}
}
