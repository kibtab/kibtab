package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoRoot is the path from this package to the repository root.
const repoRoot = "../.."

// TestDeployFiles checks each setting that a release depends on.
// It needs no container, so the unit run covers the deploy files.
func TestDeployFiles(t *testing.T) {
	cases := []struct {
		name  string
		file  string
		wants []string
	}{
		{
			name: "the instance image",
			file: "Dockerfile",
			wants: []string{
				"FROM scratch",
				"ENV CGO_ENABLED=0",
				"COPY --from=kibtab-build /out/kibtab /kibtab",
				"USER 65532:65532",
				`ENTRYPOINT ["/kibtab"]`,
				"target kibtab",
			},
		},
		{
			name: "the proxy image",
			file: "Dockerfile",
			wants: []string{
				"FROM caddy:2-builder",
				"--with github.com/mholt/caddy-ratelimit",
				"FROM caddy:2-alpine",
			},
		},
		{
			name: "the proxy settings",
			file: "Caddyfile",
			wants: []string{
				"{$KIBTAB_DOMAIN:localhost}",
				"rate_limit",
				"events 120",
				"window 1m",
				"reverse_proxy kibtab:8080",
			},
		},
		{
			name: "the stack",
			file: "docker-compose.yml",
			wants: []string{
				"\n  caddy:\n",
				"\n  kibtab:\n",
				"\n  db:\n",
				"KIBTAB_TABLES:",
				`KIBTAB_HOST: "0.0.0.0"`,
				"postgres:16-alpine",
				`"127.0.0.1:8080:8080"`,
				"condition: service_healthy",
			},
		},
		{
			name: "the build context",
			file: ".dockerignore",
			wants: []string{
				"client/excel/node_modules",
				"client/excel/dist",
				".git",
				".env",
			},
		},
		{
			name: "the release build",
			file: ".goreleaser.yaml",
			wants: []string{
				"CGO_ENABLED=0",
				"tar.gz",
				"zip",
				"checksums.txt",
				"docs/THIRD_PARTY_NOTICES.md",
			},
		},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			body := readRepoFile(t, testCase.file)
			for _, want := range testCase.wants {
				if !strings.Contains(body, want) {
					t.Errorf("%s: %q is missing", testCase.file, want)
				}
			}
		})
	}
}

// readRepoFile reads one file of the repository root.
func readRepoFile(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(repoRoot, name)
	body, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(body)
}
