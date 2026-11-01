//go:build e2e

// The end-to-end tests check the two official install paths. Each test
// skips when a tool that it needs is missing.
//
// Run them all in the repository root:
//
//	go test -tags e2e ./cmd/kibtab/ -v

package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// e2eUser is the user of each write. The audit row holds it.
const e2eUser = "e2e"

// docker runs one docker command in the root of the repository.
func docker(t *testing.T, args ...string) string {
	t.Helper()
	output, err := dockerTry(args...)
	if err != nil {
		t.Fatalf("docker %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	return output
}

// dockerTry runs one docker command and it returns the output and the
// error of the command.
func dockerTry(args ...string) (string, error) {
	cmd := exec.Command("docker", args...)
	cmd.Dir = repoRoot
	cmd.Env = os.Environ()
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	err := cmd.Run()
	return output.String(), err
}

// freePort returns a port that no process holds.
func freePort(t *testing.T) int {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("find a free port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("close the probe listener: %v", err)
	}
	return port
}

// probe sends one request. It returns the status and the body, and it
// reports a connection that the instance has not opened yet.
func probe(client *http.Client, base, method, path,
	body string) (int, string, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, base+path, reader)
	if err != nil {
		return 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Kibtab-User", e2eUser)
	response, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}
	payload, readErr := io.ReadAll(response.Body)
	closeErr := response.Body.Close()
	if closeErr != nil {
		return 0, "", closeErr
	}
	if readErr != nil {
		return 0, "", readErr
	}
	return response.StatusCode, string(payload), nil
}

// request sends one request and it fails the test when nothing answers.
func request(t *testing.T, client *http.Client, base, method, path,
	body string) (int, string) {
	t.Helper()
	status, payload, err := probe(client, base, method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, base+path, err)
	}
	return status, payload
}

// plainClient returns the client for a route without TLS.
func plainClient() *http.Client {
	return &http.Client{Timeout: 10 * time.Second}
}

// tlsClient returns the client that accepts the local certificate of
// Caddy. The stack signs the certificate for the test domain itself.
func tlsClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			//nolint:gosec // the stack holds a local certificate only.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// waitProcess waits for a process that the test stopped. It logs an
// unexpected exit only. A process that the test killed reports the code
// -1, so the kill is not an error of the test.
func waitProcess(t *testing.T, name string, cmd *exec.Cmd) {
	t.Helper()
	err := cmd.Wait()
	if err == nil {
		return
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) && exit.ExitCode() == -1 {
		return
	}
	t.Logf("wait for the %s: %v", name, err)
}

// waitHealthy waits until base answers the health route. The start of
// Caddy closes a connection, so the wait repeats each error.
func waitHealthy(t *testing.T, client *http.Client, base string) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		status, _, err := probe(client, base, http.MethodGet, "/healthz", "")
		if err == nil && status == http.StatusOK {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("%s did not answer the health route", base)
}

// step is one request of the sync walk with the value that it must hold.
type step struct {
	name   string
	method string
	path   string
	body   string
	status int
	holds  []string
}

// syncSteps returns each request of the sync walk for one table.
func syncSteps(table string) []step {
	write := `{"version":1,"key":"A1",` +
		`"cells":[{"field":"name","value":"widget"}]}`
	rows := "/v1/tables/" + table + "/rows"
	return []step{
		{
			name:   "the table list holds the table",
			method: http.MethodGet,
			path:   "/v1/tables",
			status: http.StatusOK,
			holds:  []string{`"name":"` + table + `"`},
		},
		{
			name:   "the write returns the new version",
			method: http.MethodPost,
			path:   rows,
			body:   write,
			status: http.StatusOK,
			holds:  []string{`"version":2`},
		},
		{
			name:   "the row list holds the value",
			method: http.MethodGet,
			path:   rows,
			status: http.StatusOK,
			holds:  []string{`"value":"widget"`},
		},
		{
			name:   "the audit row holds the user",
			method: http.MethodGet,
			path:   rows + "/A1/audit?field=name",
			status: http.StatusOK,
			holds:  []string{`"user":"` + e2eUser + `"`, `"new_value":"widget"`},
		},
		{
			name:   "the stale write returns the conflict",
			method: http.MethodPost,
			path:   rows,
			body:   write,
			status: http.StatusConflict,
			holds:  []string{`"conflict"`},
		},
	}
}

// walkSync walks the sync route on base. It checks each status and each
// value that the body must hold.
func walkSync(t *testing.T, client *http.Client, base, table string) {
	t.Helper()
	for _, step := range syncSteps(table) {
		status, body := request(t, client, base, step.method, step.path,
			step.body)
		if status != step.status {
			t.Errorf("%s: got status %d, want %d, body %s", step.name,
				status, step.status, body)
			continue
		}
		for _, want := range step.holds {
			if !strings.Contains(body, want) {
				t.Errorf("%s: %q is missing from %s", step.name, want, body)
			}
		}
	}
}
