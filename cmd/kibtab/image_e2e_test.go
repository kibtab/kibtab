//go:build e2e

// The image test builds the instance image and it runs that image
// against a PostgreSQL container. It checks the container itself: the
// base, the user, and the entrypoint.

package main

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"
)

const (
	// imageTag is the tag of the instance image under test.
	imageTag = "kibtab-e2e-image"
	// imageNetwork is the network that joins the image and the database.
	imageNetwork = "kibtab-e2e-network"
	// imageDB is the name of the database container.
	imageDB = "kibtab-e2e-db"
	// imageApp is the name of the instance container.
	imageApp = "kibtab-e2e-app"
	// imageTable is the table that the image test writes.
	imageTable = "image_check"
)

// removeImageStack removes each container and the network of the test.
func removeImageStack(t *testing.T) {
	t.Helper()
	for _, name := range []string{imageApp, imageDB} {
		if output, err := dockerTry("rm", "-f", name); err != nil {
			t.Logf("remove %s: %v\n%s", name, err, output)
		}
	}
	if output, err := dockerTry("network", "rm", imageNetwork); err != nil {
		t.Logf("remove the network: %v\n%s", err, output)
	}
}

// waitImageDatabase waits until the database container accepts a query.
func waitImageDatabase(t *testing.T) {
	t.Helper()
	deadline := time.Now().Add(60 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := dockerTry("exec", imageDB, "pg_isready", "-U",
			"kibtab", "-d", "kibtab"); err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatal("the database container did not accept a query")
}

// checkImageMetadata checks the base, the user, and the entrypoint of the
// image. A scratch base holds no shell, so a shell command must fail.
func checkImageMetadata(t *testing.T) {
	t.Helper()
	if output, err := dockerTry("run", "--rm", "--entrypoint", "/bin/sh",
		imageTag, "-c", "true"); err == nil {
		t.Fatalf("the image holds a shell: %s", output)
	}

	user := strings.TrimSpace(docker(t, "image", "inspect",
		"--format", "{{.Config.User}}", imageTag))
	if user != "65532:65532" {
		t.Fatalf("the image user: got %q, want 65532:65532", user)
	}

	entry := strings.TrimSpace(docker(t, "image", "inspect",
		"--format", `{{join .Config.Entrypoint " "}}`, imageTag))
	if entry != "/kibtab" {
		t.Fatalf("the entrypoint: got %q, want /kibtab", entry)
	}
}

// TestImageEndToEnd builds the instance image and it runs that image
// with a PostgreSQL container. It then walks the sync route on the
// published port of the container.
func TestImageEndToEnd(t *testing.T) {
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker is not installed")
	}
	t.Cleanup(func() { removeImageStack(t) })

	docker(t, "build", "--target", "kibtab", "-t", imageTag, ".")
	checkImageMetadata(t)

	docker(t, "network", "create", imageNetwork)
	docker(t, "run", "-d", "--name", imageDB, "--network", imageNetwork,
		"-e", "POSTGRES_USER=kibtab",
		"-e", "POSTGRES_PASSWORD=image-check",
		"-e", "POSTGRES_DB=kibtab",
		"postgres:16-alpine")
	waitImageDatabase(t)

	port := freePort(t)
	docker(t, "run", "-d", "--name", imageApp, "--network", imageNetwork,
		"-p", fmt.Sprintf("127.0.0.1:%d:8080", port),
		"-e", "PORT=8080",
		"-e", "KIBTAB_HOST=0.0.0.0",
		"-e", "DATABASE_URL=postgres://kibtab:image-check@"+imageDB+
			":5432/kibtab?sslmode=disable",
		"-e", "KIBTAB_TABLES="+imageTable,
		imageTag)

	base := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := plainClient()
	waitHealthy(t, client, base)
	docker(t, "exec", imageDB, "psql", "-U", "kibtab", "-d", "kibtab",
		"-c", "CREATE TABLE IF NOT EXISTS "+imageTable+
			" (k text PRIMARY KEY, name text)")

	walkSync(t, client, base, imageTable)
}
