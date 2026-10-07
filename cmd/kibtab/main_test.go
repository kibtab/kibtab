package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"
)

func TestHealthzAnswersWithStatusAndVersion(t *testing.T) {
	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthzHandler(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status code: got %d, want %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content type: got %q, want application/json", got)
	}
	var body healthResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode body %q: %v", recorder.Body.String(), err)
	}
	if body.Status != "ok" {
		t.Fatalf("status field: got %q, want %q", body.Status, "ok")
	}
	if body.Version != version {
		t.Fatalf("version field: got %q, want %q", body.Version, version)
	}
}

func TestHealthzRejectsNonGetMethod(t *testing.T) {
	methods := []string{
		http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch,
	}
	for _, method := range methods {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(method, "/healthz", nil)
		healthzHandler(recorder, request)
		if recorder.Code != http.StatusMethodNotAllowed {
			t.Fatalf("%s: got %d, want %d", method, recorder.Code,
				http.StatusMethodNotAllowed)
		}
		if got := recorder.Header().Get("Allow"); got != http.MethodGet {
			t.Fatalf("%s: Allow header: got %q, want %q", method, got,
				http.MethodGet)
		}
	}
}

func TestListenAddress(t *testing.T) {
	cases := []struct {
		name    string
		env     map[string]string
		want    string
		wantErr bool
	}{
		{name: "default", env: map[string]string{}, want: "127.0.0.1:8080"},
		{name: "set", env: map[string]string{"PORT": "9001"},
			want: "127.0.0.1:9001"},
		{name: "one", env: map[string]string{"PORT": "1"},
			want: "127.0.0.1:1"},
		{name: "top", env: map[string]string{"PORT": "65535"},
			want: "127.0.0.1:65535"},
		{name: "not a number", env: map[string]string{"PORT": "http"},
			wantErr: true},
		{name: "zero", env: map[string]string{"PORT": "0"}, wantErr: true},
		{name: "negative", env: map[string]string{"PORT": "-1"}, wantErr: true},
		{name: "above the range", env: map[string]string{"PORT": "65536"},
			wantErr: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			getenv := func(key string) string { return testCase.env[key] }
			got, err := listenAddress(getenv)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("PORT=%q: got address %q, want error",
						testCase.env["PORT"], got)
				}
				return
			}
			if err != nil {
				t.Fatalf("PORT=%q: unexpected error: %v",
					testCase.env["PORT"], err)
			}
			if got != testCase.want {
				t.Fatalf("PORT=%q: got %q, want %q",
					testCase.env["PORT"], got, testCase.want)
			}
		})
	}
}

func TestErrorMessage(t *testing.T) {
	portError := errInvalidPort{port: 70000}
	if got, want := portError.Error(), "PORT must be a TCP port number, got 70000"; got != want {
		t.Fatalf("invalid port: got %q, want %q", got, want)
	}
	databaseError := errMissingDatabaseURL{}
	if got, want := databaseError.Error(), "DATABASE_URL is not set"; got != want {
		t.Fatalf("missing database URL: got %q, want %q", got, want)
	}
}

func TestHealthzLogsWriteError(t *testing.T) {
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer)
	t.Cleanup(func() { log.SetOutput(os.Stderr) })

	brokenBody := failingWriter{header: http.Header{}}
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthzHandler(brokenBody, request)

	if got := logBuffer.String(); !strings.Contains(got, "healthz: write response") {
		t.Fatalf("log: got %q, want a write failure record", got)
	}
}

// failingWriter is a response writer that fails every write. It forces the
// error path of the health encoder.
type failingWriter struct {
	header http.Header
}

func (w failingWriter) Header() http.Header { return w.header }

func (w failingWriter) WriteHeader(int) {}

func (w failingWriter) Write([]byte) (int, error) {
	return 0, io.ErrClosedPipe
}

func TestNewMuxServesHealthz(t *testing.T) {
	server := httptest.NewServer(newMux())
	t.Cleanup(server.Close)

	response, err := http.Get(server.URL + "/healthz")
	if err != nil {
		t.Fatalf("get /healthz: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status code: got %d, want %d", response.StatusCode, http.StatusOK)
	}
	var body healthResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "ok" || body.Version != version {
		t.Fatalf("body: got %+v, want status ok and version %q", body, version)
	}
}

func TestRunServesHealthCheck(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve a port: %v", err)
	}
	address := listener.Addr().String()
	port := strings.TrimPrefix(address, "127.0.0.1:")
	listener.Close()

	getenv := func(key string) string {
		if key == "PORT" {
			return port
		}
		return ""
	}

	// The serve seam records the bound listener. The test closes the
	// listener at the end, so run returns and the port frees for the next
	// run of the suite.
	ready := make(chan struct{})
	var bound net.Listener
	original := serve
	serve = func(l net.Listener, h http.Handler) error {
		bound = l
		close(ready)
		return original(l, h)
	}
	t.Cleanup(func() { serve = original })

	var logBuffer bytes.Buffer
	done := make(chan error, 1)
	go func() { done <- run(getenv, &logBuffer) }()

	select {
	case <-ready:
	case <-time.After(2 * time.Second):
		t.Fatal("the instance never reached the serve call")
	}

	response, err := http.Get("http://" + address + "/healthz")
	if err != nil {
		t.Fatalf("get /healthz from the bound port: %v", err)
	}
	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		t.Fatalf("status code: got %d, want %d", response.StatusCode, http.StatusOK)
	}
	var body healthResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		response.Body.Close()
		t.Fatalf("decode body: %v", err)
	}
	response.Body.Close()
	if body.Status != "ok" || body.Version != version {
		t.Fatalf("body: got %+v, want status ok and version %q", body, version)
	}
	if got := logBuffer.String(); !strings.Contains(got, "listens on") {
		t.Fatalf("log: got %q, want a listen record", got)
	}

	if err := bound.Close(); err != nil {
		t.Fatalf("close the listener: %v", err)
	}
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("a closed listener: got no error, want one")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("run never returned")
	}
}

func TestRunRejectsBadConfig(t *testing.T) {
	getenv := func(key string) string {
		if key == "PORT" {
			return "http"
		}
		return ""
	}
	err := run(getenv, io.Discard)
	if err == nil {
		t.Fatal("PORT=http: got no error, want one")
	}
}

func TestRunRejectsBindConflict(t *testing.T) {
	held, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("hold a port: %v", err)
	}
	defer held.Close()
	address := held.Addr().String()
	port := strings.TrimPrefix(address, "127.0.0.1:")

	getenv := func(key string) string {
		if key == "PORT" {
			return port
		}
		return ""
	}
	err = run(getenv, io.Discard)
	if err == nil {
		t.Fatal("a held port: got no error, want a bind conflict")
	}
}

// TestMainProvesTheWiringExitsOnAServeFailure. The test replaces the exit
// call, so the process stays alive. The bad PORT makes run return an error
// at once.
func TestRunReturnsTheServeError(t *testing.T) {
	failure := errors.New("serve stopped")
	serve = func(net.Listener, http.Handler) error { return failure }
	t.Cleanup(func() { serve = http.Serve })

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve a port: %v", err)
	}
	port := strings.TrimPrefix(listener.Addr().String(), "127.0.0.1:")
	listener.Close()

	getenv := func(key string) string {
		if key == "PORT" {
			return port
		}
		return ""
	}
	err = run(getenv, io.Discard)
	if err != failure {
		t.Fatalf("serve failure: got %v, want %v", err, failure)
	}
}

func TestRunReturnsNoErrorAfterAServedSession(t *testing.T) {
	serve = func(net.Listener, http.Handler) error { return nil }
	t.Cleanup(func() { serve = http.Serve })

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve a port: %v", err)
	}
	port := strings.TrimPrefix(listener.Addr().String(), "127.0.0.1:")
	listener.Close()

	getenv := func(key string) string {
		if key == "PORT" {
			return port
		}
		return ""
	}
	err = run(getenv, io.Discard)
	if err != nil {
		t.Fatalf("a served session: got %v, want no error", err)
	}
}

func TestMainProvesTheWiringExitsOnAServeFailure(t *testing.T) {
	original := exit
	code := 0
	exit = func(value int) { code = value }
	t.Cleanup(func() { exit = original })

	t.Setenv("PORT", "http")
	main()

	if code != 1 {
		t.Fatalf("exit code: got %d, want 1", code)
	}
}

func TestCheckDatabaseURL(t *testing.T) {
	cases := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{name: "set", value: "postgres://localhost/kibtab"},
		{name: "empty", value: "", wantErr: true},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			getenv := func(string) string { return testCase.value }
			err := checkDatabaseURL(getenv)
			if testCase.wantErr && err == nil {
				t.Fatal("empty DATABASE_URL: got no error, want one")
			}
			if !testCase.wantErr && err != nil {
				t.Fatalf("set DATABASE_URL: unexpected error: %v", err)
			}
		})
	}
}
