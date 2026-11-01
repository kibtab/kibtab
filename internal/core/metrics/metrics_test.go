// Package metrics holds the metrics for the Kibtab instance.
// Prometheus text format metrics are served at GET /metrics.
package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	if m == nil {
		t.Fatal("NewMetrics should not return nil")
	}
}

func TestIncrRequestsTotal(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsTotal()
	m.IncrRequestsTotal()
	m.IncrRequestsTotal()
	if m.requestsTotal != 3 {
		t.Fatalf("expected requestsTotal = 3, got %d", m.requestsTotal)
	}
}

func TestTrackMetricsDoesNotPanic(t *testing.T) {
	m := NewMetrics()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := TrackMetrics(handler, m)
	if wrapped == nil {
		t.Fatal("TrackMetrics should not return nil")
	}
}

func TestWriteMetricsReturnsOK(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsTotal()

	recorder := httptest.NewRecorder()
	http.NewRequest(http.MethodGet, "/", nil) // just to verify it compiles; request not used

	m.WriteMetrics(recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected StatusOK, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !contains(body, "kibtab_requests_total") {
		t.Fatal("expected kibtab_requests_total in response body")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestObserveRequestLatencyDoesNotPanic(t *testing.T) {
	m := NewMetrics()
	d := 50 * time.Millisecond
	m.ObserveRequestLatency(d)
	// Should not panic
}