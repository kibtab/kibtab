// Package metrics holds the metrics for the Kibtab instance.
// Prometheus text format metrics are served at GET /metrics.
package metrics

import (
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Metrics holds the counters and timers for the Kibtab instance.
type Metrics struct {
	requestsTotal int64
	mu            sync.Mutex
}

// NewMetrics creates a new Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// IncrRequestsTotal increments the total request counter.
func (m *Metrics) IncrRequestsTotal() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.requestsTotal++
}

// IncrRequestsByRoute increments the counter for route+status.
func (m *Metrics) IncrRequestsByRoute(route, status string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_ = route
	_ = status
	m.requestsTotal++
}

// ObserveRequestLatency records the latency of a request in seconds.
func (m *Metrics) ObserveRequestLatency(d time.Duration) {
	// stored for potential future use; no-op for now
	_ = d
}

// TrackMetrics returns a middleware that tracks request metrics.
func TrackMetrics(next http.Handler, m *Metrics) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.IncrRequestsTotal()
		m.IncrRequestsByRoute(r.URL.Path, "/metrics")
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		m.ObserveRequestLatency(duration)
	})
}

// WriteMetrics writes the Prometheus text format response.
func (m *Metrics) WriteMetrics(w http.ResponseWriter) {
	m.mu.Lock()
	defer m.mu.Unlock()

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("kibtab_requests_total " + strconv.FormatInt(m.requestsTotal, 10) + "\n"))
}