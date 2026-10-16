package http

import (
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// failingWriter fails on every write. It covers the write error of encodeJSON.
type failingWriter struct {
	header nethttp.Header
}

func (w *failingWriter) Header() nethttp.Header { return w.header }

func (w *failingWriter) WriteHeader(int) {}

func (w *failingWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestEncodeJSON(t *testing.T) {
	recorder := httptest.NewRecorder()

	err := encodeJSON(recorder, nethttp.StatusOK, ErrorResponse{
		Code:    "ok",
		Message: "done",
	})
	if err != nil {
		t.Fatalf("encodeJSON: %v", err)
	}
	if recorder.Code != nethttp.StatusOK {
		t.Fatalf("status: got %d, want %d", recorder.Code, nethttp.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != contentTypeJSON {
		t.Fatalf("content type: got %q, want %q", got, contentTypeJSON)
	}
	if got, want := recorder.Body.String(), `{"code":"ok","message":"done"}`; got != want {
		t.Fatalf("body: got %q, want %q", got, want)
	}
}

func TestEncodeJSONMarshalError(t *testing.T) {
	recorder := httptest.NewRecorder()

	err := encodeJSON(recorder, nethttp.StatusOK, make(chan int))
	if err == nil {
		t.Fatal("encodeJSON: want an error for an unsupported value")
	}
}

func TestEncodeJSONWriteError(t *testing.T) {
	writer := &failingWriter{header: nethttp.Header{}}

	err := encodeJSON(writer, nethttp.StatusOK, ErrorResponse{})
	if err == nil {
		t.Fatal("encodeJSON: want an error when the write fails")
	}
}

func TestDecodeJSON(t *testing.T) {
	request := httptest.NewRequest(nethttp.MethodPost, "/v1/tables",
		strings.NewReader(`{"code":"c","message":"m"}`))

	var got ErrorResponse
	if err := decodeJSON(request, &got); err != nil {
		t.Fatalf("decodeJSON: %v", err)
	}
	if got.Code != "c" || got.Message != "m" {
		t.Fatalf("decodeJSON: got %+v", got)
	}
}

func TestDecodeJSONError(t *testing.T) {
	request := httptest.NewRequest(nethttp.MethodPost, "/v1/tables",
		strings.NewReader("{"))

	var got ErrorResponse
	if err := decodeJSON(request, &got); err == nil {
		t.Fatal("decodeJSON: want an error for malformed JSON")
	}
}
