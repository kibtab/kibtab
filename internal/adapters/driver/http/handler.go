package http

import (
	"log/slog"
	nethttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/kibtab/kibtab/internal/core/ports"
)

// handlers holds the ports and the logger of the HTTP transport.
type handlers struct {
	tables      ports.TableRegistry
	rows        ports.RowRepository
	syncService ports.SyncService
	log         *slog.Logger
}

// NewHandler builds the HTTP routes of the transport.
func NewHandler(tables ports.TableRegistry, rows ports.RowRepository,
	syncService ports.SyncService, log *slog.Logger) nethttp.Handler {
	h := handlers{
		tables:      tables,
		rows:        rows,
		syncService: syncService,
		log:         log,
	}

	router := chi.NewRouter()
	router.Get("/v1/tables", h.listTables)
	router.Get("/v1/tables/{table}/rows", h.listRows)
	router.Post("/v1/tables/{table}/rows", h.addRow)
	return router
}

// writeJSON writes value as a JSON response. It logs a write failure.
func (h handlers) writeJSON(w nethttp.ResponseWriter, status int, value any) {
	if err := encodeJSON(w, status, value); err != nil {
		h.log.Error("write the response", "error", err)
	}
}
