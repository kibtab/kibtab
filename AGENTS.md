# Agent Rules for Kibtab

This document gives instructions for AI coding agents.
You must follow these rules when you write code or edit files in this repository.

## 1. Code Architecture Rules

### 1.1 Layer Rules

* Keep code in the correct Hexagonal layer.
* `internal/core/` must contain only pure Go domain logic and ports.
* Do not import a database driver into `internal/core/`.
* Do not import an HTTP package into `internal/core/`.
* Define each interface in `internal/core/ports/`. Implement it in an adapter.
* Put each database engine in its own package under
  `internal/adapters/driven/<engine>/`.
* Put each transport in its own package under
  `internal/adapters/driver/<transport>/`.
* Put all CLI commands inside `cmd/kibtab/`. The CLI holds no domain logic.

The canonical list of ports is in `plan.md` section 2.2.
Add a port there first. Then define it in `internal/core/ports/`.

### 1.2 Database Agnostic Rules

Kibtab is a server for any relational database.
The core must not name one engine.

* Put every SQL statement inside the engine package that owns it.
* Do not put a PostgreSQL type, function, or DSN in `internal/core/`.
* Reach a database only through the ports in `internal/core/ports/`.
* Use the `Dialect` port for quoting, the type map, and the paging.
* Give each engine its own `Dialect` value. Do not switch on the engine name in
  the core.
* Keep a table name and a column name as data. Validate them against the
  registry before a query uses them.
* Add a new database by adding one package under `driven/` and a constructor
  in the wiring. Do not change a file under `internal/core/`.

### 1.3 Client Agnostic Rules

Kibtab is a client for any spreadsheet.
The core must not name one client.

* Treat the `SyncService` port as the boundary. A spreadsheet talks to the
  contract, not to a handler.
* Keep the value model client-neutral. Do not model an Excel range in the core.
* Put a client-specific mapping in the transport package that serves it.
* Do not put a client name, such as Excel or Sheets, in `internal/core/`.
* Keep the sync request and the sync result in the domain package. Every
  transport reuses them.
* Add a new spreadsheet by adding a client adapter and a manifest. Do not
  change a file under `internal/core/`.
* Add a new transport by adding one package under `driver/`. Do not change a
  file under `internal/core/`.

### 1.4 Wiring Rules

* Build each adapter in one place. Put that code under `cmd/kibtab/`.
* Read the engine choice and the transport choice from the environment.
* Fail at start when a chosen adapter has no implementation.
* Keep `internal/core/` free of build tags. A build tag changes an adapter,
  never the domain.

## 2. Writing Code

* Use Go 1.22 or a newer version.
* Write clear Go code. Use `gofmt` to format all code.
* Handle all errors directly. Do not ignore returned errors.
* Use `pgx/v5` for the PostgreSQL engine. Add another driver for another
  engine.
* Keep functions short. Do not write functions with more than 50 lines.
* Keep files short. Do not write files with more than 300 lines.
* Write table-driven tests for each service in `internal/core/services/`.
* Run `CGO_ENABLED=0 go build ./...`, `CGO_ENABLED=0 go vet ./...`, and
  `CGO_ENABLED=0 go test ./...` before you commit.

## 3. Writing Documentation

Write all technical documentation in two styles at the same time.
The first style is Simplified Technical English (STE100).
The second style is British English.
Both rules apply to every document. Neither style replaces the other.

### 3.1 Simplified Technical English

* Load the vendored skill before you write documentation. Read
  `./skills/simple-english/SKILL.md` first.
* Keep sentences shorter than 20 words.
* Use active voice.
* Use approved words from `./skills/simple-english/references/word-swaps.md`.
* Use one term for one meaning in each document.
* Put the condition before the command.
* Define each technical term at its first use.

### 3.2 British English

* Use British English spellings in every document, comment, and commit message.
* Use `-ise`, not `-ize`. The forms are `organise`, `recognise`, `analyse`.
* Use `-our`, not `-or`. The forms are `colour`, `behaviour`, `favour`.
* Use `-re`, not `-er`. The forms are `centre`, `metre`, `litre`.
* Use `licence` as a noun. Use `license` as a verb.
* Keep the proper names of licences. The name is `Apache License 2.0`.
* Double the final `l` after a vowel. The forms are `travelling`, `cancelled`.
* Use `artefact`, not `artifact`.
* Use `whilst` and `amongst` only where the sense needs them.
* Keep the plural `-ae` and `-a` forms. The form is `formulae`.

### 3.3 Rules For Both Styles

* Do not use em dashes, contractions, or the words `should`, `would`,
  `may`, and `might` in technical documentation.
* Keep articles. Keep `that` after a verb.
* Put the version in the first line of every document.
* Keep copied work as it was written upstream. Read the rule below.
* Run `make docs-check` before you commit. The target checks the sentence
  length and the banned words. It does not check the word lists.

### 3.4 Copied Work

Some files hold text from another project.
Keep that text as it was written upstream.
Do not apply the rules in sections 3.1 to 3.3 to it.

| File | Source |
| --- | --- |
| `skills/simple-english/` | The SimpleEnglish skill. MIT licence. |
| `CODE_OF_CONDUCT.md` | Contributor Covenant 3.0. CC BY-SA 4.0. |

A contributor report and a pull request are not documentation.
Do not apply the rules in section 3 to them.

## 4. Database Rules

Write a rule here for every engine. Never write one rule for all engines.
The rules below hold for the engine they name.

### 4.1 Rules For Every Engine

* Do not execute raw SQL string concatenations.
* Use parameterized queries to prevent SQL injection.
* Wrap a multi-cell update in one transaction of the target engine.
* Quote each identifier through the `Dialect` port. Do not quote it by hand.
* Never accept a table or column name from a request without validating it
  against the table registry.
* Keep each statement in the package for its engine.

### 4.2 Rules For The PostgreSQL Engine

* Use `pgx/v5` for the connection.
* Quote each identifier with `pgx.Identifier`.
* Wrap a multi-cell update with `BEGIN` and `COMMIT`.
* Run the migrations under a PostgreSQL advisory lock.

## 5. Build And Toolchain Rules

* Set `CGO_ENABLED=0` in every build, test, and release command.
* Do not add a dependency that needs CGO. Do not use `go-sqlite3`,
  `net` C resolvers, or DNS libraries in C.
* Keep `go.mod` free of `require` blocks that pull a C compiler.
* Run the engine locally with `air`. Read `air.toml` before you change it.
* Keep the local stack in `docker-compose.yml`.
* Do not keep the ASD-STE100 word lists in this repository. Use the vendored
  skill in `./skills/simple-english/` for the Simplified Technical English
  rules.

## 6. Release Rules

* Release with GoReleaser. Read `.goreleaser.yaml` before you change it.
* Use semantic versioning. Follow `plan.md` for the version meaning.
* Write one file at `docs/changelogs/vX.Y.Z.md` for each release.
* Follow the format in `docs/changelogs/README.md` for each file.
* Copy `docs/changelogs/TEMPLATE.md` before you write the file.
* Update `docs/index.md` and `docs/CREDITS.md` in the same commit as the
  release notes.
* Update `docs/THIRD_PARTY_NOTICES.md` when a dependency changes.
* Tag the commit. Do not move a published tag.

## 7. Repository Rules

* Do not commit secrets, `.env` files, or credentials.
* Do not commit build output, the `dist/` folder, or coverage files.
* Keep the module path as `github.com/kibtab/kibtab`.
* Put the project plan in `plan.md`. Update it when the plan changes.
* Follow `plan.md` for scope. Do not build a feature that a later version
  holds.
* Every task in `plan.md` uses a markdown check box. Mark it done when the
  code and the release note exist. Delete the line when the release ships.

## 8. Commit Message Rules

Kibtab uses Conventional Commits.
Write each message for a reader who knows the plan but not the change.

### 8.1 Format

```text
<type>(<scope>): <description>
```

Keep the subject under 50 characters.
Use the imperative mood. Start with a verb.
Do not end the subject with a full stop.
Use British English in the subject and in the body.

### 8.2 Types

| Type | Use for |
| --- | --- |
| `feat` | A new feature for a user. |
| `fix` | A defect repair. |
| `docs` | A documentation change only. |
| `refactor` | A code change with no new behaviour. |
| `test` | A test change only. |
| `perf` | A change that improves a measured result. |
| `build` | The toolchain, the dependency, or the build settings. |
| `ci` | The workflow files. |
| `chore` | A task with no source change. |

### 8.3 Scopes

Use a scope from the layer or the area of the change.
A change to more than one layer uses the widest scope.

| Scope | Covers |
| --- | --- |
| `core` | `internal/core/`. The domain, the ports, the services. |
| `adapters` | The wiring of each adapter in `cmd/kibtab/`. |
| `postgres` | `internal/adapters/driven/postgres/`. |
| `duckdb` | `internal/adapters/driven/duckdb/`. |
| `http` | `internal/adapters/driver/http/`. The routes and the codec. |
| `cli` | `cmd/kibtab/`. The entry point and the flags. |
| `client` | `client/`. The spreadsheet clients. |
| `schema` | The `_kibtab_meta` schema and the migrations. |
| `deploy` | The `Dockerfile`, the `Caddyfile`, and the compose file. |
| `release` | GoReleaser, the tags, and the archives. |
| `docs` | Every file under `docs/`, plus `README.md` and `AGENTS.md`. |
| `deps` | `go.mod`, `go.sum`, and the client packages. |
| `ci` | The workflow files. |
| `skills` | The vendored skill in `skills/simple-english/`. |

An engine scope names the engine. Use `duckdb` for that adapter.
Use `core` for a change in the domain. It is not an engine scope.

### 8.4 Examples

```text
feat(core): add the optimistic version compare
fix(postgres): stop the audit insert from rolling back the version bump
docs(core): record the port rules in the architecture guide
feat(http): add the POST /v1/tables/{table}/cells route
refactor(postgres): move the row query out of the table registry
feat(duckdb): add the DuckDB dialect behind the Dialect port
build(release): set CGO_ENABLED=0 in the goreleaser environment
test(postgres): cover the rollback of a failed batch
chore(skills): pin the simple-english skill to 32ea2d3
```

### 8.5 Scope Selection Rules

* Choose the scope of the file that changed most.
* Use `deps` for a version bump. Use `release` for a build setting.
* Use `core` for a change in more than one package under `internal/core/`.
* Use an engine scope for a change in one engine package.
* Use `adapters` for a change in the wiring that builds each adapter.
* Do not invent a scope. Add a new scope to the table in this section first.
* Use no scope when the change touches the whole repository.

### 8.6 Breaking Changes

A change to a REST path, a JSON field name, or a CLI flag is breaking.
Write `!` after the scope. Then add a `BREAKING CHANGE:` footer.

```text
feat(http)!: replace the cell write body with a versioned envelope

The v1 field names move to the envelope meta object.
A client must read the version from meta.version.

BREAKING CHANGE: the v1 JSON field names are removed.
Read docs/changelogs/v2.0.0.md before you upgrade.
```

### 8.7 The Body

Add a body when the subject does not explain the change.
Wrap each line at 72 characters.
Explain why the change was needed. Not what the diff shows.
Link the issue at the end with `Closes #123`.

### 8.8 What An Agent Must Do After A Change

When you finish a change, suggest a commit message.
Do not run the commit yourself unless the user asks for it.

Give the message in a fenced block.
Use the format in section 8.1 with a type and a scope from section 8.3.
State why the scope fits the change in one sentence.
Name the version in `plan.md` that the change belongs to.
Say whether the change needs a MAJOR version.

If the change has no good scope, say so. Propose a new scope and add it to
the table in section 8.3.

## 9. Scope And Discipline Rules

### 9.1 The Scope Of A Change

* Do one thing in a commit. Split the rest into another commit.
* Do one thing in a function. Split the rest into another function.
* Solve the problem in the task. Do not solve the next task.
* Keep each commit to one purpose. Keep each file to one purpose.
* Reject a change that touches two layers for one reason. Split it.
* Do not add a flag, a switch, or a mode for a case that no user has.
* Do not add an abstraction before a second caller needs it.
* Keep the diff small. A large diff hides a large defect.

### 9.2 The Unix Rules

* Do one thing and do it well.
* Write each tool to read from a stream and to write to a stream where it
  makes sense.
* Keep each module small. Push the detail down to the layer that owns it.
* Use the standard library first. Add a dependency only for a real need.
* Fail loudly and early. Do not hide an error behind a default value.
* Prefer text a person can read over text only a tool can parse.
* Keep the core free of a convenience for any one adapter.

### 9.3 The RTFM Rules

* Read the file before you change it. Read the whole file.
* Read `plan.md` before you start. Read `AGENTS.md` before you edit.
* Search for the answer before you ask. Then search for the question.
* Point to the document that holds the answer. Do not repeat it in a reply.
* Write the answer in one file. Link to that file from every other place.
* Do not copy a setup step into a second file. One file owns each step.
* When a rule has no answer in the repository, write the rule in the
  repository first. Then follow it.
* Update the document in the same commit as the code it describes.

### 9.4 The Scope Test

Run this test before you send a change for review.

1. Name the one thing the change does. Stop when you need a second verb.
2. Name the layer that owns the change. Stop when you need a second layer.
3. Name the version in `plan.md` that holds the change.
4. Read the diff. Remove each hunk that serves a second purpose.
5. Check that the commit message names that one thing.

## 10. Separation Of Concerns Rules

Keep each part of the code for one reason only.
Give each file, each package, and each guide one job.

### 10.1 The Layer Rules

* Keep each file to one job. Split a file that serves two jobs.
* Keep each package to one job. Split a package that serves two jobs.
* Keep the domain free of transport, persistence, and presentation.
* Keep the transport free of SQL. Keep SQL free of the domain.
* Keep the wiring in `cmd/kibtab/`. Do not build an adapter elsewhere.
* Keep a configuration value out of the domain. Read it at the edge.

### 10.2 The Change Rules

* Change one concern in a commit. Split the rest into another commit.
* Do not mix a repair with a new feature in one commit.
* Do not mix a rename with a behaviour change in one commit.
* Move code before you change it. Use one commit for the move.
* Do not add a second reason to keep an argument in a function.
* Return a value. Do not set a field for a caller in another layer.

### 10.3 The Documentation Rules

* Give each guide one task. Do not copy a step into a second file.
* Put each setup step in the one file that owns it.
* Put a section of many files in its own folder. Add a `README.md` index.
* Link each file of a section from the index. Do not leave a file unlinked.
* Update the index in the same commit as the file that it lists.
* Read the index before you write a new file. Place the file beside its task.

### 10.4 The Test Rules

* Test one concern per test. Split a test that asserts two things.
* Name the test for the concern that it covers.
* Keep the core test free of an adapter. Use a fake for a port.
* Run the contract suite for an adapter. Do not repeat it in each adapter.
* Keep a fixture in one file. Share it. Do not copy it into each test.