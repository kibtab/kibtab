package http

import (
	"io"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

// discardLogger returns a logger that writes nowhere.
func discardLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestNewHandlerUnknownRoute(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeAudit{}, &fakeSync{},
		discardLogger())

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/nope", nil))

	if recorder.Code != nethttp.StatusNotFound {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusNotFound)
	}
}
