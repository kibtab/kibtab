# Kibtab Implementation Plan

Version 0.0.0-plan. This document holds the plan for Kibtab.
It splits the work into releases.
Each release has one version number.
The version number follows semantic versioning.

Read the sections in order.
Do not start work on a later release before the earlier release works.

## How To Use This Plan

Every task in this plan is a markdown check box.
A check box holds work that is not done.

* Mark a box with `- [x]` when the code and the release note exist.
* Delete a line when the release ships.
* Never mark a box done before the code compiles and the test passes.

The plan holds no status field.
The check box holds the status.

Each release needs a file at `docs/changelogs/vX.Y.Z.md`.
The format is in `docs/changelogs/README.md`.
The template is in `docs/changelogs/TEMPLATE.md`.
No release file exists yet.

## 1. Version Meaning

Kibtab uses the version number to tell a user what changed.

| Part | Meaning in Kibtab |
| --- | --- |
| MAJOR | The public interface changed. A client must change. |
| MINOR | Kibtab added a feature. An old client keeps working. |
| PATCH | Kibtab fixed a defect. No interface changed. |

Kibtab holds the REST paths, the JSON field names, and the CLI flags as the
public interface.
A change to any of them needs a MAJOR version.

| Version | Theme | Result |
| --- | --- | --- |
| v0.1.0 | Skeleton | The engine starts. It answers a health check. |
| v0.2.0 | Domain and ports | The pure Go logic exists. It has tests. |
| v0.3.0 | PostgreSQL adapter | Kibtab reads and writes rows. |
| v0.4.0 | Write path | A cell change writes in a transaction. |
| v0.5.0 | HTTP API | A client reads and writes over REST. |
| v0.6.0 | Office.js client | A user edits a cell in Excel. |
| v0.7.0 | Audit and locking | Kibtab records each change. It blocks a stale write. |
| v0.8.0 | Packaging | GoReleaser ships binaries. Air runs the engine. |
| v0.9.0 | Hardening | The engine survives load. It reports its health. |
| v1.0.0 | Stable | The interface freezes. |
| v1.1.0 and later | Growth | New adapters and new features. |

## 2. Proposed Codebase Structure

Kibtab must stay flexible.
It must serve any relational database.
It must serve any spreadsheet.
The structure below keeps both sides replaceable.

### 2.1 The Layout

```text
kibtab/
├── cmd/
│   └── kibtab/                 # The wiring. It builds each adapter.
│       └── main.go
├── internal/
│   ├── core/                   # The kernel. It has no I/O.
│   │   ├── domain/             # The models. They name no engine or client.
│   │   ├── ports/              # The interfaces. The core owns them.
│   │   └── services/           # The use cases. They call the ports only.
│   └── adapters/
│       ├── driven/             # One package per database engine.
│       │   └── postgres/       # The first engine.
│       └── driver/             # One package per transport or client.
│           └── http/           # The first transport.
├── client/                     # The spreadsheet clients. One folder each.
├── docs/                       # The documentation and the changelogs.
├── skills/simple-english/      # The vendored writing skill.
├── AGENTS.md
├── plan.md
├── air.toml
├── Caddyfile
├── docker-compose.yml
├── Dockerfile
├── .goreleaser.yaml
└── go.mod
```

### 2.2 The Port Names

The core names no engine and no client.
It names a need.
The adapter supplies the answer.

| Port | Need it fills |
| --- | --- |
| `RowRepository` | Read and write rows in a table. |
| `TableRegistry` | List the tables. Read the metadata for a table. |
| `AuditWriter` | Record a change to a cell. |
| `Dialect` | Quote an identifier. Map a value to a type. Page a query. |
| `TransactionRunner` | Run a group of writes as one unit. |
| `Clock` | Give the current time. |
| `SyncService` | Accept a sync payload. Return a sync result. |

A new engine implements the ports. It does not change the core.
A new client implements `SyncService`. It does not change the core.

### 2.3 The Rules For Agility

* Keep the core free of an engine name.
* Keep the core free of a client name.
* Keep the core free of a transport name.
* Give each engine its own package under `adapters/driven/`.
* Give each transport its own package under `adapters/driver/`.
* Use `Clock` in the core. Do not call the system clock in a test.
* Reach a database only through a port.
* Put the wiring in `cmd/kibtab/`. Keep it in one place.
* Choose the engine and the transport from the environment.
* Add a CI step. It fails the build when the core imports an adapter.

### 2.4 The Test Strategy

Agility needs a test that holds each side to the contract.

* [ ] Write a contract test suite for `RowRepository`. Run it for each engine.
* [ ] Write a contract test suite for `SyncService`. Run it for each client.
* [ ] Run the engine suite in CI against a real database container.
* [ ] Keep the core test suite free of a container.
* [ ] Measure the engine swap. Add a second engine with no core change.

## 3. Rules For All Releases

Every release below must meet these rules.
Mark each box only when the check passes.

* [ ] `CGO_ENABLED=0 go build ./...` succeeds with no environment variable.
* [ ] The module has no dependency that needs a C compiler.
* [ ] All SQL sits in the engine package under `internal/adapters/driven/`.
* [ ] `internal/core/` has no database driver import and no HTTP import.
* [ ] `internal/core/` names no database engine and no spreadsheet client.
* [ ] `CGO_ENABLED=0 go vet ./...` reports no problem.
* [ ] `make docs-check` reports no problem.
* [ ] The documentation follows `AGENTS.md` section 3.
* [ ] Each guide owns one task. No setup step sits in two files.
* [ ] Each section of many files has an index that links each file.
* [ ] A changelog file exists at `docs/changelogs/vX.Y.Z.md`.

## 4. Release v0.1.0 - Skeleton

**Goal:** The engine starts and it answers a health check.

### Scope

* [ ] Create `go.mod`. Set the module path and Go 1.22.
* [ ] Create `cmd/kibtab/main.go`. It reads `PORT` and `DATABASE_URL`.
* [ ] Keep `README.md` as a landing page. Link to each guide under `docs/`.
* [ ] Keep `CONTRIBUTING.md`. It holds the steps to send a change.
* [ ] Keep `CODE_OF_CONDUCT.md`. It holds the terms for taking part.
* [ ] Keep `SECURITY.md`. It holds the private route for a vulnerability.
* [ ] Keep `.github/ISSUE_TEMPLATE/`. It holds the forms for an issue.
* [ ] Add the CI step. It fails the build when a module has no licence row
      in `docs/THIRD_PARTY_NOTICES.md`.
* [ ] Keep `docs/install.md`. It holds the install and build steps.
* [ ] Keep `docs/licence.md`. It holds the terms for the engine and the skill.
* [ ] Create the folder `internal/core/` with `domain`, `ports`, `services`.
* [ ] Create the folder `internal/adapters/driven/` for the engines.
* [ ] Create the folder `internal/adapters/driver/` for the transports.
* [ ] Create the folder `client/` for the spreadsheet clients.
* [ ] Add the route `GET /healthz`. It returns the status and the version.
* [ ] Add `air.toml`. Air rebuilds on a change in `internal/` or `cmd/`.
* [ ] Add the `Makefile` with the targets `build`, `test`, `lint`, `run`,
      `docs-check`, `notice-check`, and `release-check`.
* [ ] Add `scripts/docs-check.py`. It checks the sentence length and the
      banned words. It does not need the ASD word lists.
* [ ] Add `scripts/notice-check.py`. It compares `go.mod` with the tables in
      `docs/THIRD_PARTY_NOTICES.md`.
* [ ] Add `.github/workflows/ci.yml`. It builds, vets, tests, and checks the
      documentation.
* [ ] Set `CGO_ENABLED=0` in every step of the workflow.

### Out Of Scope

This release holds no database access.
This release holds no write path.

### Done When

* [ ] `CGO_ENABLED=0 go build ./...` succeeds on Linux, macOS, and Windows.
* [ ] `air` starts the engine and reloads it on a code change.
* [ ] `curl localhost:8080/healthz` returns the version.
* [ ] `docs/changelogs/v0.1.0.md` exists.

## 5. Release v0.2.0 - Domain And Ports

**Goal:** The pure domain exists and it has tests.

### Scope

* [ ] Add `internal/core/domain/`. It holds `TableMetadata`, `CellValue`,
      `CellDelta`, `SyncPayload`, `SyncResult`, and `ValidationError`.
* [ ] Add `internal/core/ports/`. It holds `RowRepository`, `TableRegistry`,
      `AuditWriter`, `Dialect`, `TransactionRunner`, `Clock`, and
      `SyncService`.
* [ ] Add `internal/core/services/`. It holds the validation, the batch
      merge, and the version compare.
* [ ] Write a table-driven test for each service.
* [ ] Keep `docs/architecture.md`. It holds the layers and the port contract.
* [ ] Add a CI step. It fails the build on a driver import in the core.
* [ ] Add a CI step. It fails the build when the core names an engine or a
      client.

### Out Of Scope

This release holds no PostgreSQL.
This release holds no HTTP.

### Done When

* [ ] `CGO_ENABLED=0 go test ./internal/core/...` passes with no database.
* [ ] The CI step finds no driver import in `internal/core/`.
* [ ] `docs/changelogs/v0.2.0.md` exists.

## 6. Release v0.3.0 - PostgreSQL Adapter

**Goal:** Kibtab reads and writes rows in PostgreSQL.

### Scope

* [ ] Add `internal/adapters/driven/postgres/`. It uses `pgx/v5`.
* [ ] Create the schema `_kibtab_meta`.
* [ ] Create the table `_kibtab_meta.table_versions`.
* [ ] Create the table `_kibtab_meta.audit_log`.
* [ ] Add the migrations. They run at start. They use an advisory lock.
* [ ] Implement the `Dialect` port for PostgreSQL. It quotes with
      `pgx.Identifier`. It maps each type. It pages a query.
* [ ] Implement `TableRegistry`. It reads the table list from the
      environment.
* [ ] Implement `RowRepository`. It reads a page of rows. It writes one row.
* [ ] Implement `TransactionRunner`. It opens a transaction.
* [ ] Use a parameterized query for each statement.
* [ ] Add an integration test. It writes a row and reads it back.
* [ ] Add the row for `pgx/v5` to `docs/THIRD_PARTY_NOTICES.md`.

### Out Of Scope

This release holds no cell write path.

### Done When

* [ ] The integration test writes a row and reads it back.
* [ ] A second start runs the migrations again with no error.
* [ ] No query in the adapter uses a string concatenation for a value.
* [ ] `make notice-check` passes.
* [ ] `docs/changelogs/v0.3.0.md` exists.

## 7. Release v0.4.0 - Write Path

**Goal:** A cell change writes to the database in a transaction.

### Scope

* [ ] Add the multi-cell write service. It groups the changes by row.
* [ ] Wrap each group in a transaction. Use `BEGIN` and `COMMIT`.
* [ ] Add the optimistic lock. Compare the version before the write.
* [ ] Add the value coercion. It maps a string to the column type.
* [ ] Return a conflict when the version does not match.
* [ ] Add a race test. Two changes use one version. One write wins.
* [ ] Add a rollback test. A failed group does not stop the other groups.

### Out Of Scope

This release holds no audit log.
This release holds no REST API.

### Done When

* [ ] The race test passes. One write wins.
* [ ] The rollback test passes. The other groups stay.
* [ ] `docs/changelogs/v0.4.0.md` exists.

## 8. Release v0.5.0 - HTTP API

**Goal:** A client can read and write over REST.

### Scope

* [ ] Add the handlers in `internal/adapters/driver/http/`.
* [ ] Add the route `GET /v1/tables`.
* [ ] Add the route `GET /v1/tables/{table}/rows`.
* [ ] Add the route `POST /v1/tables/{table}/rows`.
* [ ] Add the route `POST /v1/tables/{table}/cells`.
* [ ] Add the JSON codec. It lives beside the handlers.
* [ ] Add the request size limit. Reject a body above the limit.
* [ ] Add the CORS policy for the Office.js origin.
* [ ] Add the structured log. Write one line per request.
* [ ] Add the error type. It holds a code and a message.
* [ ] Add `docs/openapi.yaml`. It holds the REST interface definition.
* [ ] Add a contract test for each route.
* [ ] Add the rows for the router and the UUID library to
      `docs/THIRD_PARTY_NOTICES.md`.

### Out Of Scope

This release holds no Office.js taskpane.

### Done When

* [ ] Each route has a passing contract test.
* [ ] The OpenAPI file matches each handler.
* [ ] `docs/changelogs/v0.5.0.md` exists.

## 9. Release v0.6.0 - Office.js Client

**Goal:** A user edits a cell in Excel and Kibtab stores it.

### Scope

* [ ] Create `client/`. It holds the Office.js taskpane.
* [ ] The taskpane sends a `SyncPayload` and shows the result.
* [ ] The taskpane shows a conflict message. It offers a reload.
* [ ] The taskpane groups a burst of edits into one payload.
* [ ] Add `manifest.xml` for the add-in.
* [ ] Add the client build. It uses `esbuild`.
* [ ] Add the CI job for the client. It runs apart from the Go build.
* [ ] Add the output `kibtab-taskpane.zip`.
* [ ] Add the rows for the client packages to
      `docs/THIRD_PARTY_NOTICES.md`.

### Out Of Scope

This release holds no Google Sheets code.
This release holds no store submission.

### Done When

* [ ] A manual test edits a cell in Excel. `psql` shows the new value.
* [ ] The build produces the `kibtab-taskpane.zip` file.
* [ ] `docs/changelogs/v0.6.0.md` exists.

## 10. Release v0.7.0 - Audit And Locking

**Goal:** Kibtab records each change and it blocks a stale write.

### Scope

* [ ] Write each cell change to `_kibtab_meta.audit_log`.
* [ ] Store the user, the time, the old value, and the new value.
* [ ] Bump the version. Run it in the same transaction as the data write.
* [ ] Add the query path. It reads the audit rows for one cell.
* [ ] Add the retention job. It removes old audit rows.
* [ ] Read the retention days from `KIBTAB_AUDIT_RETENTION_DAYS`.
* [ ] Add a test. It writes a cell and finds the audit row.

### Out Of Scope

This release holds no export of an audit row.

### Done When

* [ ] The test writes a cell and finds the audit row.
* [ ] A stale write returns the conflict code. It writes nothing.
* [ ] `docs/changelogs/v0.7.0.md` exists.

## 11. Release v0.8.0 - Packaging

**Goal:** GoReleaser ships binaries. Air runs the engine in development.

### Scope

* [ ] Add `.goreleaser.yaml`. It builds for Linux, macOS, and Windows.
* [ ] Set `CGO_ENABLED=0` in the GoReleaser build environment.
* [ ] Build an archive for each target. Use `tar.gz` for Linux and macOS.
      Use `zip` for Windows.
* [ ] Add the checksum file for each release.
* [ ] Read the release notes from `docs/changelogs/`.
* [ ] Add `docs/THIRD_PARTY_NOTICES.md` to each archive.
* [ ] Keep `air.toml` at the repository root as the only reload file.
* [ ] Add `Dockerfile`. Use a multi-stage build and a `scratch` base.
* [ ] Add `docker-compose.yml`. It holds Caddy, the engine, and PostgreSQL.
* [ ] Add `Caddyfile`. It holds the domain and the rate limit.
* [ ] Add the target `make release-check`. It runs `goreleaser check`.
* [ ] Add the rows for GoReleaser, Caddy, and PostgreSQL to
      `docs/THIRD_PARTY_NOTICES.md`.
* [ ] Add the section `docs/self-hosting/`. It holds the guides for a server.
* [ ] Add the guide `docs/self-hosting/README.md`. It indexes the section.
* [ ] Add the guide `docs/self-hosting/deployment.md`. It owns the first start.
* [ ] Add the guide `docs/self-hosting/configuration.md`. It owns each setting.
* [ ] Add the guide `docs/self-hosting/upgrade.md`. It owns each version change.
* [ ] Add the guide `docs/self-hosting/backup-restore.md`. It owns the data.
* [ ] Add the guide `docs/self-hosting/troubleshooting.md`. It owns each fault.
* [ ] Link each self-hosting guide from the section index and from
      `docs/index.md`.

### Out Of Scope

The default job does not push an image to a registry.
This section holds no Kubernetes guide.

### Done When

* [ ] `goreleaser check` passes.
* [ ] A dry run produces an archive for each target.
* [ ] The image starts and answers the health check.
* [ ] `ldd` reports no dynamic library for the Linux binary.
* [ ] `docs/changelogs/v0.8.0.md` exists.

## 12. Release v0.9.0 - Hardening

**Goal:** The engine survives load and it reports its health.

### Scope

* [ ] Add the connection pool settings. Read them from the environment.
* [ ] Add the rate limit. Caddy sets the limit at the edge.
* [ ] Add the route `GET /metrics`. It uses the Prometheus text format.
* [ ] Set a timeout on each query.
* [ ] Add the metric labels. They are the route, the status, and the time.
* [ ] Add the load test. It uses `k6`.
* [ ] Add `docs/runbook.md`.
* [ ] Add the row for the metrics library to
      `docs/THIRD_PARTY_NOTICES.md`.

### Out Of Scope

This release holds no horizontal scale.

### Done When

* [ ] The load test shows no error at the target rate.
* [ ] The metrics route shows the query count and the error count.
* [ ] `docs/changelogs/v0.9.0.md` exists.

## 13. Release v1.0.0 - Stable

**Goal:** The public interface freezes.

### Scope

* [ ] Freeze the REST paths and the JSON field names.
* [ ] Freeze the CLI flags and the exit codes.
* [ ] Write the migration guide for each earlier version.
* [ ] Confirm the Apache License 2.0 covers every file that Kibtab wrote.
      The licence does not change at this version.
* [ ] Sign the release artefact.
* [ ] Publish the announcement.
* [ ] Add the contract test suite. It runs against the frozen interface.

### Out Of Scope

This release adds no feature.

### Done When

* [ ] The contract test suite passes against the frozen interface.
* [ ] `docs/changelogs/v1.0.0.md` lists each change that breaks an interface
      from v0.x.
* [ ] Each migration guide exists.

## 14. Releases After v1.0.0

These releases are the expected direction.
The boxes stay open until the work starts.

### v1.1.0 - Google Sheets Adapter

* [ ] Add a second inbound adapter. It uses Apps Script.
* [ ] Send the same REST payload from the adapter.

### v1.2.0 - Google Sheets Add-On

* [ ] Publish the add-on to the Google Workspace Marketplace.
* [ ] Move the taskpane code to a shared package.

### v1.3.0 - Schema Discovery

* [ ] Read the database schema. Build the table registry from it.
* [ ] Drop the need for a table list in the environment.

### v1.4.0 - Value Types

* [ ] Support the date type, the boolean type, the decimal type, and the
      enum type.
* [ ] Validate a value against the column type before the write.

### v1.5.0 - Live Sync

* [ ] Add a WebSocket channel.
* [ ] Push each change to each open taskpane.

### v1.6.0 - Offline Queue

* [ ] Store a change in the taskpane when the network fails.
* [ ] Send the change when the network returns.

### v1.7.0 - Row Security

* [ ] Read the user claims in the engine.
* [ ] Add a row filter to each query.

### v1.8.0 - Multi-Database

* [ ] Add a driven adapter for DuckDB.
* [ ] Keep the core package free of a change.

### v1.9.0 - Query Views

* [ ] Store a named filter with a table.
* [ ] Show the view in the taskpane.

### v2.0.0 - New Public Interface

* [ ] Replace the JSON body of the write route with a versioned envelope.
* [ ] Remove the field names from v1.
* [ ] Publish a migration guide for each v1 route.

### v3.0.0 - Plugin Ports

* [ ] Open the port interfaces to an external module.
* [ ] Load a module at start.
* [ ] Let a module add a driven adapter.

## 15. Risk And Mitigation

| Risk | Effect | Mitigation |
| --- | --- | --- |
| A write race loses data | High | The version check runs in the transaction. |
| A large table slows the page | Medium | The page size has a limit. The key has an index. |
| Caddy holds a certificate rate limit | Medium | The stack supports a test domain first. |
| A query appears in the core layer | High | CI fails the build on a driver import. |
| A dependency needs CGO | High | The build sets `CGO_ENABLED=0`. The build then fails. |
| The interface changes after v1.0.0 | High | The contract test suite freezes the interface. |
| An upstream skill change alters the docs | Medium | The skill is vendored at a pinned commit. |

## 16. Definition Of Done For A Release

A release is done when all of these are true.

* The code builds with `CGO_ENABLED=0` for each target.
* The unit tests and the integration tests pass.
* `CGO_ENABLED=0 go vet ./...` reports no problem.
* `make docs-check` reports no problem.
* The changelog file exists at `docs/changelogs/vX.Y.Z.md`.
* `docs/index.md` lists the release.
* The changelog states the upgrade steps.
* The tag exists. The tag does not move.
