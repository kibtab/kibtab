package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
