//go:build e2e

// The release test runs the GoReleaser binary on the host with a native
// PostgreSQL and a native Caddy. It checks the second install path: the
// release archive on a bare host behind a proxy.

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const (
	// releaseTable is the table that the release test writes.
	releaseTable = "release_check"
	// releaseDB is the name of the native database.
	releaseDB = "kibtab"
	// releaseUser is the superuser of the native database.
	releaseUser = "kibtab"
)

// toolFor skips the test when one tool is missing.
func toolFor(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s is not installed", name)
	}
}

// runTool runs one command and it returns the output. It fails the test
// when the command fails.
func runTool(t *testing.T, dir string, env []string, name string,
	args ...string) string {
	t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), env...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return string(output)
}

// buildReleaseBinary builds the release binary with GoReleaser and it
// returns the path of the binary.
func buildReleaseBinary(t *testing.T) string {
	t.Helper()
	runTool(t, repoRoot, nil, "goreleaser", "build", "--snapshot", "--clean",
		"--single-target")
	matches, err := filepath.Glob(
		filepath.Join("..", "..", "dist", "*", "kibtab"))
	if err != nil {
		t.Fatalf("find the release binary: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("the release build wrote no binary")
	}
	binary, err := filepath.Abs(matches[0])
	if err != nil {
		t.Fatalf("resolve the release binary: %v", err)
	}
	return binary
}

// startPostgres starts a native database in the work folder. It returns
// the connection string and the port of the database.
func startPostgres(t *testing.T, work string) (string, int) {
	t.Helper()
	data := filepath.Join(work, "data")
	socket := filepath.Join(work, "socket")
	if err := os.MkdirAll(socket, 0o700); err != nil {
		t.Fatalf("make the socket folder: %v", err)
	}
	runTool(t, work, nil, "initdb", "-D", data, "-U", releaseUser,
		"--auth-local=trust", "--auth-host=trust", "--no-locale")
	port := freePort(t)
	runTool(t, work, nil, "pg_ctl", "-D", data,
		"-l", filepath.Join(work, "postgres.log"),
		"-o", fmt.Sprintf("-k %s -p %d -h 127.0.0.1", socket, port),
		"start")
	t.Cleanup(func() {
		command := exec.Command("pg_ctl", "-D", data, "stop", "-m", "fast")
		if err := command.Run(); err != nil {
			t.Logf("stop the database: %v", err)
		}
	})
	runTool(t, work, nil, "createdb", "-h", "127.0.0.1",
		"-p", fmt.Sprint(port), "-U", releaseUser, releaseDB)
	return fmt.Sprintf("postgres://%s@127.0.0.1:%d/%s?sslmode=disable",
		releaseUser, port, releaseDB), port
}

// createReleaseTable makes the table that the test writes.
func createReleaseTable(t *testing.T, work string, port int) {
	t.Helper()
	runTool(t, work, nil, "psql", "-h", "127.0.0.1", "-p",
		fmt.Sprint(port), "-U", releaseUser, "-d", releaseDB, "-c",
		"CREATE TABLE IF NOT EXISTS "+releaseTable+
			" (k text PRIMARY KEY, name text)")
}

// startInstance runs the release binary and it returns the port of the
// instance.
func startInstance(t *testing.T, binary, databaseURL string) int {
	t.Helper()
	port := freePort(t)
	cmd := exec.Command(binary)
	cmd.Dir = repoRoot
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("PORT=%d", port),
		"KIBTAB_TABLES="+releaseTable,
		"DATABASE_URL="+databaseURL,
	)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the instance: %v", err)
	}
	t.Cleanup(func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("stop the instance: %v", err)
		}
		waitProcess(t, "instance", cmd)
	})
	return port
}

// startProxy runs Caddy in front of the instance and it returns the port
// of the proxy.
func startProxy(t *testing.T, work string, instancePort int) int {
	t.Helper()
	port := freePort(t)
	config := filepath.Join(work, "Caddyfile")
	body := fmt.Sprintf("{\n\tadmin off\n}\n\nhttp://localhost:%d {\n"+
		"\treverse_proxy 127.0.0.1:%d\n}\n", port, instancePort)
	if err := os.WriteFile(config, []byte(body), 0o600); err != nil {
		t.Fatalf("write the proxy config: %v", err)
	}
	cmd := exec.Command("caddy", "run", "--config", config,
		"--adapter", "caddyfile")
	cmd.Dir = work
	cmd.Env = append(os.Environ(),
		"XDG_DATA_HOME="+filepath.Join(work, "caddy-data"),
		"XDG_CONFIG_HOME="+filepath.Join(work, "caddy-config"),
	)
	if err := cmd.Start(); err != nil {
		t.Fatalf("start the proxy: %v", err)
	}
	t.Cleanup(func() {
		if err := cmd.Process.Kill(); err != nil {
			t.Logf("stop the proxy: %v", err)
		}
		waitProcess(t, "proxy", cmd)
	})
	return port
}

// TestReleaseEndToEnd runs the release binary with a native PostgreSQL
// and a native Caddy. It walks the sync route through the proxy.
func TestReleaseEndToEnd(t *testing.T) {
	for _, name := range []string{"goreleaser", "initdb", "pg_ctl", "psql",
		"createdb", "caddy"} {
		toolFor(t, name)
	}
	work := t.TempDir()
	binary := buildReleaseBinary(t)
	databaseURL, databasePort := startPostgres(t, work)
	createReleaseTable(t, work, databasePort)

	instancePort := startInstance(t, binary, databaseURL)
	client := plainClient()
	waitHealthy(t, client, fmt.Sprintf("http://127.0.0.1:%d", instancePort))

	proxyPort := startProxy(t, work, instancePort)
	proxyBase := fmt.Sprintf("http://localhost:%d", proxyPort)
	waitHealthy(t, client, proxyBase)

	walkSync(t, client, proxyBase, releaseTable)
}
