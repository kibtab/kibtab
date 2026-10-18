package http

import (
	"bytes"
	"log/slog"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// bufferLogger returns a logger that writes one line to the buffer.
func bufferLogger(buffer *bytes.Buffer) *slog.Logger {
	return slog.New(slog.NewTextHandler(buffer, nil))
}

func TestRequestLog(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeSync{},
		bufferLogger(&buffer))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))

	if got := strings.Count(buffer.String(), "\n"); got != 1 {
		t.Fatalf("lines: got %d, want 1, body %q", got, buffer.String())
	}
	for _, want := range []string{
		"msg=request",
		"method=GET",
		"path=/v1/tables",
		"status=200",
		"duration_ms=",
		"id=",
	} {
		if !strings.Contains(buffer.String(), want) {
			t.Fatalf("log line %q misses %q", buffer.String(), want)
		}
	}
	if recorder.Header().Get(requestIDHeader) == "" {
		t.Fatalf("%s: got no identifier", requestIDHeader)
	}
}

func TestRequestLogNotFound(t *testing.T) {
	var buffer bytes.Buffer
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeSync{},
		bufferLogger(&buffer))
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/nope", nil))

	if !strings.Contains(buffer.String(), "status=404") {
		t.Fatalf("log line %q misses the 404 status", buffer.String())
	}
}
