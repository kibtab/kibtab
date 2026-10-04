# Kibtab Makefile
# The default goal prints the help text.

BINARY := bin/kibtab
GO ?= go
GORELEASER ?= goreleaser
COVERAGE_FLOOR := 80
CHANGELOG_URL_BASE := https://github.com/kibtab/kibtab/blob/main/docs/changelogs

.DEFAULT_GOAL := help

.PHONY: help build install run test cover cover-html cover-verify coverage-svg \
        vet lint fmt tidy watch release-test release-check snapshot \
        bench bench-all bench-save bench-compare \
        profile-cpu profile-trace \
        docs-check notice-check check \
        tag clean tools

help: ## Show this help
	@printf "kibtab (spreadsheet as client, database as server)\n\n"
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-14s\033[0m %s\n", $$1, $$2}'

build: ## Build the kibtab binary into bin/
	CGO_ENABLED=0 $(GO) build -trimpath -ldflags "-s -w" -o $(BINARY) ./cmd/kibtab

install: ## Install the CLI into $GOPATH/bin (go install)
	CGO_ENABLED=0 $(GO) install -trimpath -ldflags "-s -w" ./cmd/kibtab

run: ## Run the CLI (extra args after --)
	CGO_ENABLED=0 $(GO) run ./cmd/kibtab -- $(filter-out $@,$(MAKECMDGOALS))

test: ## Run all tests with a coverage summary
	CGO_ENABLED=0 $(GO) test -cover ./...

cover: ## Run the tests and print the per-function coverage breakdown
	CGO_ENABLED=0 $(GO) test -coverpkg=./... -coverprofile=coverage.out ./... -count=1
	@$(GO) tool cover -func=coverage.out

coverage-svg: ## Regenerate the coverage badge from the test run
	CGO_ENABLED=0 $(GO) test -coverpkg=./... -coverprofile=coverage.out ./... -count=1
	@go-test-coverage -p coverage.out -b coverage.svg 2>/dev/null || \
		echo "install go-test-coverage to regenerate the badge (make tools)"

cover-html: ## Run the tests and open the HTML coverage report in a browser
	CGO_ENABLED=0 $(GO) test -coverpkg=./... -coverprofile=coverage.out ./... -count=1
	$(GO) tool cover -html=coverage.out

cover-verify: ## Fail when module coverage sits below the floor
	CGO_ENABLED=0 $(GO) test -coverpkg=./... -coverprofile=coverage.out ./... -count=1
	@total=$$($(GO) tool cover -func=coverage.out | awk '/^total:/ {gsub("%",""); print $$3}'); \
	echo "total coverage: $$total% (floor $(COVERAGE_FLOOR)%)"; \
	rm -f coverage.out; \
	if [ "$$(echo "$$total < $(COVERAGE_FLOOR)" | bc -l)" = "1" ]; then \
		echo "coverage below the $(COVERAGE_FLOOR)% floor" >&2; exit 1; \
	fi

vet: ## Run go vet over all packages
	CGO_ENABLED=0 $(GO) vet ./...

lint: ## Run the Go linter over all packages
	golangci-lint run ./...

fmt: ## Format all Go source with gofmt
	gofmt -w .

tidy: ## Tidy the Go module files
	$(GO) mod tidy

watch: ## Hot-reload cmd/kibtab on save (needs air)
	air

docs-check: ## Check the documentation for sentence length and banned words
	python3 scripts/docs-check.py plan.md AGENTS.md README.md docs CONTRIBUTING.md SECURITY.md

notice-check: ## Compare go.mod with the dependency table in the notices
	python3 scripts/notice-check.py

check: build vet test docs-check notice-check ## Run every gate before a commit

release-test: ## Dry-run the release: build every target into dist/ (no upload)
	$(MAKE) build
	$(MAKE) snapshot

release-check: ## Validate the GoReleaser configuration
	CGO_ENABLED=0 $(GORELEASER) check

snapshot: ## Test GoReleaser locally in snapshot mode
	CGO_ENABLED=0 $(GORELEASER) release --snapshot --clean

bench: ## Run the benchmarks of the core services
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. -benchmem ./internal/core/...

bench-all: ## Run the benchmarks of every package that carries one
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. -benchmem ./... -count=1

bench-save: ## Save a benchmark run as the comparison baseline (bench.txt)
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. -benchmem ./internal/core/... -count=1 | tee bench.txt

bench-compare: ## Compare a fresh benchmark run with bench.txt (needs benchstat)
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. -benchmem ./internal/core/... -count=6 > new-bench.txt
	@benchstat bench.txt new-bench.txt 2>/dev/null || \
		echo "install benchstat to read the comparison (make tools)"

profile-cpu: ## Write a CPU profile of the benchmark run into cpu.prof
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. -benchmem ./internal/core/... -cpuprofile=cpu.prof

profile-trace: ## Write an execution trace of the benchmark run into trace.out
	CGO_ENABLED=0 $(GO) test -run=^$$ -bench=. ./internal/core/... -trace=trace.out

tag: ## Tag the suggested version (the highest changelog) and push the tag
	@HIGHEST=$$(ls docs/changelogs/v*.md 2>/dev/null | sed -E 's|.*/v([0-9]+\.[0-9]+\.[0-9]+)\.md|\1|' | sort -V | tail -1); \
	if [ -z "$$HIGHEST" ]; then \
		CURRENT=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
		MAJOR=$$(echo "$$CURRENT" | sed 's/^v//' | cut -d. -f1); \
		MINOR=$$(echo "$$CURRENT" | sed 's/^v//' | cut -d. -f2); \
		HIGHEST="$$MAJOR.$$(($$MINOR + 1)).0"; \
	fi; \
	SUGGEST="v$$HIGHEST"; \
	if git rev-parse "$$SUGGEST" >/dev/null 2>&1; then \
		MAJOR=$$(echo "$$HIGHEST" | cut -d. -f1); \
		MINOR=$$(echo "$$HIGHEST" | cut -d. -f2); \
		SUGGEST="v$$MAJOR.$$(($$MINOR + 1)).0"; \
	fi; \
	read -p "Enter version [$$SUGGEST]: " TAG; \
	TAG=$${TAG:-$$SUGGEST}; \
	NOTES="$(CHANGELOG_URL_BASE)/$$TAG.md"; \
	if [ ! -f "docs/changelogs/$$TAG.md" ]; then \
		echo "Warning: docs/changelogs/$$TAG.md is missing. Write the release notes before the release."; \
	fi; \
	if git rev-parse "$$TAG" >/dev/null 2>&1; then \
		echo "Tag $$TAG already exists, pushing..."; \
	else \
		git tag -a "$$TAG" -m "Release $$TAG. Notes: $$NOTES" && echo "Created tag $$TAG with the notes link."; \
	fi; \
	git push origin "$$TAG"

clean: ## Remove the build artefacts
	CGO_ENABLED=0 $(GO) clean -cache -test-cache 2>/dev/null || true
	rm -rf bin dist coverage.out build-errors.log bench.txt new-bench.txt \
		cpu.prof trace.out *.test

tools: ## Install the development tools (air, goreleaser, golangci-lint)
	go install github.com/air-verse/air@latest
	go install github.com/goreleaser/goreleaser/v2@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/vladopajic/go-test-coverage/v2@latest
	go install golang.org/x/perf/cmd/benchstat@latest