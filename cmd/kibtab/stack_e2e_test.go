//go:build e2e

// The end-to-end test starts the compose stack and it drives the API
// through Caddy. It needs Docker and Docker Compose on the host.
//
// Run it with:
//
//	go test -tags e2e ./cmd/kibtab/ -run TestStackEndToEnd -v

package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// e2eProject is the compose project of the test stack. A separate name
// keeps the test away from the stack of a developer.
const e2eProject = "kibtab-e2e"

// e2eTable is the table that the end-to-end test writes.
const e2eTable = "stack_check"

// e2eUser is the user of the write. The audit row holds it.
const e2eUser = "e2e"

// e2eEnv is the environment of the test stack. The compose file reads
// each value from the environment of the command.
func e2eEnv() []string {
	return []string{
		"POSTGRES_PASSWORD=stack-check",
		"KIBTAB_DOMAIN=localhost",
		"KIBTAB_TABLES=" + e2eTable,
	}
}

// compose runs one docker compose command in the root of the repository.
func compose(t *testing.T, args ...string) string {
	t.Helper()
	command := append([]string{
		"compose", "-p", e2eProject, "--env-file", "/dev/null",
	}, args...)
	cmd := exec.Command("docker", command...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), e2eEnv()...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	if err := cmd.Run(); err != nil {
		t.Fatalf("docker %s: %v\n%s", strings.Join(args, " "), err,
			output.String())
	}
	return output.String()
}

// stopStack removes the containers, the network, and the data of the test.
func stopStack(t *testing.T) {
	t.Helper()
	command := []string{
		"compose", "-p", e2eProject, "--env-file", "/dev/null", "down", "-v",
	}
	cmd := exec.Command("docker", command...)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(), e2eEnv()...)
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Logf("stop the test stack: %v\n%s", err, output)
	}
}

// proxyClient returns a client that accepts the local certificate of
// Caddy. The stack signs the certificate for the test domain itself.
func proxyClient() *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			//nolint:gosec // the stack holds a local certificate only.
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
}

// waitHealthy waits until the proxy answers the health route. The start
// of Caddy closes a connection, so each try reports the error and the
// wait repeats it.
func waitHealthy(t *testing.T, client *http.Client) {
	t.Helper()
	deadline := time.Now().Add(90 * time.Second)
	for time.Now().Before(deadline) {
		status, _, err := probe(client, http.MethodGet, "/healthz", "")
		if err == nil && status == http.StatusOK {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatal("the stack did not answer the health route")
}

// probe sends one request to the proxy. It returns the status and the
// body, and it reports a connection that the stack has not opened yet.
func probe(client *http.Client, method, path, body string) (int, string, error) {
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req, err := http.NewRequest(method, "https://localhost"+path, reader)
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

// request sends one request to the proxy and it returns the status and
// the body. It fails the test when the stack does not answer.
func request(t *testing.T, client *http.Client, method, path,
	body string) (int, string) {
	t.Helper()
	status, payload, err := probe(client, method, path, body)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	return status, payload
}

// expect holds one step of the end-to-end walk.
type expect struct {
	name   string
	method string
	path   string
	body   string
	status int
	holds  []string
}

// TestStackEndToEnd starts the stack and it writes one row through the
// proxy. It reads the row and the audit row back. It then sends the same
// version again and it waits for the conflict.
func TestStackEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not installed")
	}
	t.Cleanup(func() { stopStack(t) })

	compose(t, "up", "-d", "--build")
	client := proxyClient()
	waitHealthy(t, client)
	compose(t, "exec", "-T", "db", "psql", "-U", "kibtab", "-d", "kibtab",
		"-c", "CREATE TABLE IF NOT EXISTS "+e2eTable+
			" (k text PRIMARY KEY, name text)")

	write := `{"version":1,"key":"A1",` +
		`"cells":[{"field":"name","value":"widget"}]}`
	steps := []expect{
		{
			name:   "the table list holds the table",
			method: http.MethodGet,
			path:   "/v1/tables",
			status: http.StatusOK,
			holds:  []string{`"name":"` + e2eTable + `"`},
		},
		{
			name:   "the write returns the new version",
			method: http.MethodPost,
			path:   "/v1/tables/" + e2eTable + "/rows",
			body:   write,
			status: http.StatusOK,
			holds:  []string{`"version":2`},
		},
		{
			name:   "the row list holds the value",
			method: http.MethodGet,
			path:   "/v1/tables/" + e2eTable + "/rows",
			status: http.StatusOK,
			holds:  []string{`"value":"widget"`},
		},
		{
			name:   "the audit row holds the user",
			method: http.MethodGet,
			path:   "/v1/tables/" + e2eTable + "/rows/A1/audit?field=name",
			status: http.StatusOK,
			holds:  []string{`"user":"` + e2eUser + `"`, `"new_value":"widget"`},
		},
		{
			name:   "the stale write returns the conflict",
			method: http.MethodPost,
			path:   "/v1/tables/" + e2eTable + "/rows",
			body:   write,
			status: http.StatusConflict,
			holds:  []string{`"conflict"`},
		},
	}
	for _, step := range steps {
		status, body := request(t, client, step.method, step.path, step.body)
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
