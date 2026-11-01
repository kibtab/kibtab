// Package metrics holds the metrics for the Kibtab instance.
// Prometheus text format metrics are served at GET /metrics.
package metrics

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestNewMetrics(t *testing.T) {
	m := NewMetrics()
	if m == nil {
		t.Fatal("NewMetrics should not return nil")
	}
	if m.requestsTotal != 0 {
		t.Fatalf("expected requestsTotal = 0, got %d", m.requestsTotal)
	}
}

func TestIncrRequestsTotal(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsTotal()
	if m.requestsTotal != 1 {
		t.Fatalf("expected requestsTotal = 1 after one incr, got %d", m.requestsTotal)
	}
	m.IncrRequestsTotal()
	m.IncrRequestsTotal()
	if m.requestsTotal != 3 {
		t.Fatalf("expected requestsTotal = 3 after three incrs, got %d", m.requestsTotal)
	}
}

func TestIncrRequestsByRoute(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsByRoute("/tables", "/metrics")
	if m.requestsTotal != 1 {
		t.Fatalf("expected requestsTotal = 1 after IncrRequestsByRoute, got %d", m.requestsTotal)
	}
}

func TestObserveRequestLatency(t *testing.T) {
	m := NewMetrics()
	d := 50 * time.Millisecond
	m.ObserveRequestLatency(d)
	// ObserveRequestLatency is a no-op; should not panic
}

func TestObserveRequestLatencyMultipleTimes(t *testing.T) {
	m := NewMetrics()
	m.ObserveRequestLatency(10 * time.Millisecond)
	m.ObserveRequestLatency(20 * time.Millisecond)
	m.ObserveRequestLatency(30 * time.Millisecond)
	// ObserveRequestLatency is a no-op; should not panic
}

func TestTrackMetricsDoesNotPanic(t *testing.T) {
	m := NewMetrics()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	wrapped := TrackMetrics(handler, m)
	if wrapped == nil {
		t.Fatal("TrackMetrics should not return nil")
	}
	// Verify the wrapped handler works
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	wrapped.ServeHTTP(recorder, req)
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected StatusOK from wrapped handler, got %d", recorder.Code)
	}
}

func TestWriteMetricsReturnsOK(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsTotal()

	recorder := httptest.NewRecorder()
	m.WriteMetrics(recorder)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected StatusOK, got %d", recorder.Code)
	}

	body := recorder.Body.String()
	if !contains(body, "kibtab_requests_total") {
		t.Fatal("expected kibtab_requests_total in response body")
	}
	if !contains(body, strconv.FormatInt(m.requestsTotal, 10)) {
		t.Fatal("expected request count in response body")
	}
}

func TestWriteMetricsContentType(t *testing.T) {
	m := NewMetrics()
	m.IncrRequestsTotal()

	recorder := httptest.NewRecorder()
	m.WriteMetrics(recorder)

	if recorder.Result().Header.Get("Content-Type") != "text/plain" {
		t.Fatalf("expected Content-Type text/plain, got %v", recorder.Result().Header.Get("Content-Type"))
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