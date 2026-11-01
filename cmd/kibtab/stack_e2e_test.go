//go:build e2e

// The stack test starts the compose stack and it drives the API through
// Caddy. It needs Docker and Docker Compose on the host.

package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// e2eProject is the compose project of the test stack. A separate name
// keeps the test away from the stack of a developer.
const e2eProject = "kibtab-e2e"

// e2eTable is the table that the stack test writes.
const e2eTable = "stack_check"

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

// TestStackEndToEnd starts the compose stack, writes one row through the
// proxy, and reads the row and the audit row back. It then sends the same
// version again and it waits for the conflict.
func TestStackEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not installed")
	}
	t.Cleanup(func() { stopStack(t) })

	compose(t, "up", "-d", "--build")
	client := tlsClient()
	waitHealthy(t, client, "https://localhost")
	compose(t, "exec", "-T", "db", "psql", "-U", "kibtab", "-d", "kibtab",
		"-c", "CREATE TABLE IF NOT EXISTS "+e2eTable+
			" (k text PRIMARY KEY, name text)")

	walkSync(t, client, "https://localhost", e2eTable)
}
