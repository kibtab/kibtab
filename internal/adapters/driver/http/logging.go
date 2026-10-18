package http

import (
	nethttp "net/http"
	"time"

	"github.com/google/uuid"
)

// requestIDHeader is the response header that carries the request
// identifier.
const requestIDHeader = "X-Request-Id"

// requestLog writes one structured line for each request. The line holds
// the identifier, the method, the path, the status, and the duration.
func (h handlers) requestLog(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		id := uuid.NewString()
		w.Header().Set(requestIDHeader, id)
		start := time.Now()
		recorder := &statusWriter{ResponseWriter: w, status: nethttp.StatusOK}
		next.ServeHTTP(recorder, r)
		h.log.Info("request",
			"id", id,
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
			"duration_ms", time.Since(start).Milliseconds())
	})
}

// statusWriter records the status of a response. The write passes
// through to the wrapped writer.
type statusWriter struct {
	nethttp.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}
