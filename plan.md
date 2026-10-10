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
* Never mark a box done before the code compiles and the test passes.
* Remove the checked boxes when the changelog for the version states each
  item and each deliverable exists.
* Put one line in the place of the boxes.

```markdown
Implemented, see [vX.Y.Z changelog](docs/changelogs/vX.Y.Z.md).
```

The plan holds no status field.
The check box holds the status of open work.
The replacement line holds the status of a shipped release.

Each release needs a file at `docs/changelogs/vX.Y.Z.md`.
The format is in `docs/changelogs/README.md`.
The template is in `docs/changelogs/TEMPLATE.md`.

## 1. Release Milestones

This plan holds the work, not the rules.
The rules live in [the agent rules](AGENTS.md).
The gate for a release lives in `AGENTS.md` section 5.

| Version | Theme | Result |
| --- | --- | --- |
| v0.1.0 | Skeleton | Kibtab starts. It answers a health check. |
| v0.2.0 | Domain and ports | The pure Go logic exists. It has tests. |
| v0.3.0 | PostgreSQL adapter | Kibtab reads and writes rows. |
| v0.4.0 | Write path | A cell change writes in a transaction. |
| v0.5.0 | HTTP API | A client reads and writes over REST. |
| v0.6.0 | Office.js client | A user edits a cell in Excel. |
| v0.7.0 | Audit and locking | Kibtab records each change. It blocks a stale write. |
| v0.8.0 | Packaging | GoReleaser ships binaries. Air runs the instance. |
| v0.9.0 | Hardening | The instance survives load. It reports its health. |
| v1.0.0 | Stable | The interface freezes. |
| v1.1.0 and later | Growth | New adapters and new features. |

## 2. Release v0.1.0 - Skeleton

**Goal:** Kibtab starts and it answers a health check.

### Scope

Implemented, see [v0.1.0 changelog](docs/changelogs/v0.1.0.md).

### Out Of Scope

This release holds no database access.
This release holds no write path.

## 3. Release v0.2.0 - Domain And Ports

**Goal:** The pure domain exists and it has tests.

### Scope

Implemented, see [v0.2.0 changelog](docs/changelogs/v0.2.0.md).

### Out Of Scope

This release holds no PostgreSQL.
This release holds no HTTP.

### Done When

* [x] `CGO_ENABLED=0 go test ./internal/core/...` passes with no database.
* [x] The CI step finds no driver import in `internal/core/`.
* [x] `docs/changelogs/v0.2.0.md` exists.

## 4. Release v0.3.0 - PostgreSQL Adapter

**Goal:** Kibtab reads and writes rows in PostgreSQL.

### Scope

* [x] Add `internal/adapters/driven/postgres/`. It uses `pgx/v5`.
* [x] Create the schema `_kibtab_meta`.
* [x] Create the table `_kibtab_meta.table_versions`.
* [x] Create the table `_kibtab_meta.audit_log`.
* [x] Add the migrations. They run at start. They use an advisory lock.
* [x] Implement the `Dialect` port for PostgreSQL. It quotes with
      `pgx.Identifier`. It maps each type. It pages a query.
* [x] Implement `TableRegistry`. It reads the table list from the
      environment.
* [x] Implement `RowRepository`. It reads a page of rows. It writes one row.
* [x] Implement `TransactionRunner`. It opens a transaction.
* [x] Use a parameterized query for each statement.
* [x] Write the contract test suite for `RowRepository`. Run it for each
      engine.
* [x] Run the integration suite against a real database container.
* [x] Add an integration test. It writes a row and reads it back.
* [x] Cover every error path in the adapter. Cover each boundary case.
* [x] Add the row for `pgx/v5` to `docs/THIRD_PARTY_NOTICES.md`.
* [x] Remove the `dockers` build from `.goreleaser.yaml`. The `Dockerfile`
      lands in v0.8.0.

### Out Of Scope

This release holds no cell write path.

### Done When

* [x] The integration test writes a row and reads it back.
* [x] A second start runs the migrations again with no error.
* [x] No query in the adapter uses a string concatenation for a value.
* [x] `make notice-check` passes.
* [x] `docs/changelogs/v0.3.0.md` exists.

## 5. Release v0.4.0 - Write Path

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

## 6. Release v0.5.0 - HTTP API

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
* [ ] Write the contract test suite for `SyncService`. Run it for each
      client.
* [ ] Add the rows for the router and the UUID library to
      `docs/THIRD_PARTY_NOTICES.md`.

### Out Of Scope

This release holds no Office.js taskpane.

### Done When

* [ ] Each route has a passing contract test.
* [ ] The OpenAPI file matches each handler.
* [ ] `docs/changelogs/v0.5.0.md` exists.

## 7. Release v0.6.0 - Office.js Client

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

## 8. Release v0.7.0 - Audit And Locking

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

## 9. Release v0.8.0 - Packaging

**Goal:** GoReleaser ships binaries. Air runs the instance in development.

### Scope

* [x] Add `.goreleaser.yaml`. It builds for Linux, macOS, and Windows.
* [ ] Set `CGO_ENABLED=0` in the GoReleaser build environment.
* [ ] Build an archive for each target. Use `tar.gz` for Linux and macOS.
      Use `zip` for Windows.
* [ ] Add the checksum file for each release.
* [ ] Read the release notes from `docs/changelogs/`.
* [ ] Add `docs/THIRD_PARTY_NOTICES.md` to each archive.
* [ ] Keep `air.toml` at the repository root as the only reload file.
* [ ] Set `CGO_ENABLED=0` in every GoReleaser build and in every Make target
      that calls the Go tool.
* [ ] Add `Dockerfile`. Use a multi-stage build and a `scratch` base.
* [ ] Add `docker-compose.yml`. It holds Caddy, Kibtab, and PostgreSQL.
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

## 10. Release v0.9.0 - Hardening

**Goal:** The instance survives load and it reports its health.

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

## 11. Release v1.0.0 - Stable

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

## 12. Releases After v1.0.0

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

* [ ] Read the user claims in the core.
* [ ] Add a row filter to each query.

### v1.8.0 - Multi-Database

* [ ] Add a driven adapter for DuckDB.
* [ ] Keep the core package free of a change.
* [ ] Run the `RowRepository` contract suite against DuckDB. It passes with
      no change to the suite.
* [ ] Measure the engine swap. Record the build time and the query time for
      each engine.

### v1.9.0 - Query Views

* [ ] Store a named filter with a table.
* [ ] Show the view in the taskpane.

### v2.0.0 - New Public Interface

* [ ] Replace the JSON body of the write route with a versioned envelope.
* [ ] Remove the field names from v1.
* [ ] Publish a migration guide for each v1 route.

### v3.0.0 - External Ports

* [ ] Open the port interfaces to an external module.
* [ ] Load a module at start.
* [ ] Let a module add a driven adapter.

## 13. Risk And Mitigation

| Risk | Effect | Mitigation |
| --- | --- | --- |
| A write race loses data | High | The version check runs in the transaction. |
| A large table slows the page | Medium | The page size has a limit. The key has an index. |
| Caddy holds a certificate rate limit | Medium | The stack supports a test domain first. |
| A query appears in the core layer | High | CI fails the build on a driver import. |
| A dependency needs CGO | High | The build sets `CGO_ENABLED=0`. The build then fails. |
| The interface changes after v1.0.0 | High | The contract test suite freezes the interface. |
| An upstream skill change alters the docs | Medium | The skill is vendored at a pinned commit. |

## 14. When A Release Is Done

A release is done when each check in `AGENTS.md` section 5 passes.
Read that section before you tag a release.
