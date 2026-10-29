package http

import (
	nethttp "net/http"
	"net/http/httptest"
	"testing"
)

// checkHeader asserts one response header.
func checkHeader(t *testing.T, recorder *httptest.ResponseRecorder, key, want string) {
	t.Helper()

	if got := recorder.Header().Get(key); got != want {
		t.Fatalf("%s: got %q, want %q", key, got, want)
	}
}

func TestCORSPreflight(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeAudit{}, &fakeSync{},
		discardLogger())
	request := httptest.NewRequest(nethttp.MethodOptions, "/v1/tables", nil)
	request.Header.Set("Origin", "https://localhost:3000")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != nethttp.StatusNoContent {
		t.Fatalf("status: got %d, want %d", recorder.Code,
			nethttp.StatusNoContent)
	}
	checkHeader(t, recorder, "Access-Control-Allow-Origin",
		"https://localhost:3000")
	checkHeader(t, recorder, "Access-Control-Allow-Methods", corsMethods)
	checkHeader(t, recorder, "Access-Control-Allow-Headers", corsHeaders)
}

func TestCORSRequestWithOrigin(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeAudit{}, &fakeSync{},
		discardLogger())
	request := httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil)
	request.Header.Set("Origin", "https://contoso.example")
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	checkHeader(t, recorder, "Access-Control-Allow-Origin",
		"https://contoso.example")
}

func TestCORSNoOrigin(t *testing.T) {
	handler := NewHandler(fakeRegistry{}, fakeRows{}, &fakeAudit{}, &fakeSync{},
		discardLogger())
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder,
		httptest.NewRequest(nethttp.MethodGet, "/v1/tables", nil))

	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("Access-Control-Allow-Origin: got %q, want empty", got)
	}
}
