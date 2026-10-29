package http

import (
	nethttp "net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/kibtab/kibtab/internal/core/domain"
)

// auditRowJSON is one audit row in the JSON of an audit list.
type auditRowJSON struct {
	Field     string    `json:"field"`
	User      string    `json:"user"`
	OldValue  string    `json:"old_value"`
	NewValue  string    `json:"new_value"`
	ChangedAt time.Time `json:"changed_at"`
}

// auditList is the JSON body of GET /v1/tables/{table}/rows/{row}/audit.
type auditList struct {
	Audit []auditRowJSON `json:"audit"`
}

// listAudit answers GET /v1/tables/{table}/rows/{row}/audit. The query
// parameter `field` names the cell. The route holds no other cell.
func (h handlers) listAudit(w nethttp.ResponseWriter, r *nethttp.Request) {
	table := chi.URLParam(r, "table")
	row := chi.URLParam(r, "row")
	field := r.URL.Query().Get("field")
	if field == "" {
		h.writeError(w, nethttp.StatusBadRequest, codeInvalidRequest,
			"the field query parameter is empty")
		return
	}

	entries, err := h.audit.Read(r.Context(), table, row, field)
	if err != nil {
		h.log.Error("read the audit rows", "error", err)
		h.writeError(w, nethttp.StatusInternalServerError, codeInternal,
			"the audit list failed")
		return
	}

	h.writeJSON(w, nethttp.StatusOK, auditList{Audit: toAuditRows(entries)})
}

// toAuditRows maps audit rows to the JSON of one audit list.
func toAuditRows(entries []domain.AuditRow) []auditRowJSON {
	body := make([]auditRowJSON, 0, len(entries))
	for _, entry := range entries {
		body = append(body, auditRowJSON{
			Field:     entry.Field,
			User:      entry.User,
			OldValue:  entry.OldValue,
			NewValue:  entry.NewValue,
			ChangedAt: entry.ChangedAt,
		})
	}
	return body
}
