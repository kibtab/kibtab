# Agent Rules for Kibtab

This document gives instructions for AI coding agents.
You must follow these rules when you write code or edit files in this repository.

`plan.md` holds the work. This document holds the rules.
Change this document first. Then change the code.

| Read this | For |
| --- | --- |
| [plan.md](plan.md) | The scope of each version. The open check boxes. |
| [docs/architecture.md](docs/architecture.md) | The tree, the ports, and the data flow. |
| [CONTRIBUTING.md](CONTRIBUTING.md) | The full steps to send a change. |
| [docs/licence.md](docs/licence.md) | The licence terms and the canonical documents. |
| [docs/self-hosting/](docs/self-hosting/README.md) | Running Kibtab on a server. |
| [docs/install.md](docs/install.md) | Installing and building. |

## 1. Code Architecture Rules

### 1.1 The Layers

* `internal/core/` holds pure Go domain logic and ports. It has no I/O.
* `internal/core/` imports no database driver and no HTTP package.
* Define each interface in `internal/core/ports/`. Implement it in an adapter.
* Put each database engine in its own package under
  `internal/adapters/driven/<engine>/`.
* Put each transport in its own package under
  `internal/adapters/driver/<transport>/`.
* Put each spreadsheet client in its own folder under `client/`.
* Put all CLI commands inside `cmd/kibtab/`. The CLI holds no domain logic.
* Build each adapter in one place, under `cmd/kibtab/`.
* Read the engine choice and the transport choice from the environment.
* Fail at start when a chosen adapter has no implementation.
* Keep `internal/core/` free of build tags. A build tag changes an adapter,
  never the domain.

### 1.2 The Core Names No Engine And No Client

Kibtab is a server for any relational database.
Kibtab is a client for any spreadsheet.
The core names a need. The adapter supplies the answer.

* Put every SQL statement inside the engine package that owns it.
* Do not put a PostgreSQL type, a function, or a DSN in `internal/core/`.
* Reach a database only through a port.
* Give each engine its own `Dialect` value. Never switch on the engine name.
* Keep a table name and a column name as data. Validate them against the
  registry before a query uses them.
* Treat `SyncService` as the client boundary. A spreadsheet talks to the
  contract, not to a handler.
* Keep the value model client-neutral. Do not model an Excel range in the core.
* Put no client name, such as Excel or Sheets, in `internal/core/`.
* Keep the sync request and the sync result in the domain package.
* Add a database by adding one package and a constructor. The core stays.
* Add a spreadsheet by adding a client folder and a manifest. The core stays.
* Add a transport by adding one package. The core stays.
* Add a CI step that fails the build when the core imports an adapter.

### 1.3 The Layout And The Port Contract

[docs/architecture.md](docs/architecture.md) holds the full tree, the port
table, the data flow, and the steps to add an engine, a client, or a
transport. Read it before you move a file.

Add a port to that document first.
Then define it in `internal/core/ports/`.
Use the `Clock` port in the core. Never call the system clock in a test.

## 2. Writing Code

* Use Go 1.22 or a newer version.
* Use `gofmt` on all code. Use `go vet` before you commit.
* Handle all errors directly. Do not ignore a returned error.
* Use the standard library first. Add a dependency only for a real need.
* Keep functions under 50 lines. Keep files under 300 lines.
* Write table-driven tests for each service in `internal/core/services/`.
* Return a value. Do not set a field for a caller in another layer.

## 3. Writing Documentation

Write every document in two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
Neither style replaces the other.
Both styles apply to every document, comment, and commit message.

Load the vendored skill first. Read
`./skills/simple-english/SKILL.md`.

### 3.1 Simplified Technical English

* Keep sentences shorter than 20 words.
* Use active voice and simple tenses.
* Use approved words from `./skills/simple-english/references/word-swaps.md`.
* Use one term for one meaning in each document.
* Put the condition before the command.
* Define each technical term at its first use.

### 3.2 British English

* Use `-ise`, not `-ize`: `organise`, `recognise`, `analyse`.
* Use `-our`, not `-or`: `colour`, `behaviour`, `favour`.
* Use `-re`, not `-er`: `centre`, `metre`, `litre`.
* Use `licence` as a noun. Use `license` as a verb.
* Keep the proper name of a licence: `Apache License 2.0`.
* Double the final `l` after a vowel: `travelling`, `cancelled`.
* Use `artefact`, not `artifact`.
* Use `whilst` and `amongst` only where the sense needs them.
* Keep the plural `-ae` and `-a` forms: `formulae`.

### 3.3 Rules For Both

* Do not use em dashes, contractions, or the words `should`, `would`,
  `may`, and `might`.
* Keep articles. Keep `that` after a verb.
* Put the version in the first line of every document.
* Run `make docs-check` before you commit.

### 3.4 The Canonical Document Exception

Three documents have one valid form. Copy the text as written.
[docs/licence.md](docs/licence.md) lists them and their canonical source.

Do not apply section 3 to these documents.
Do not fix a sentence, a long clause, or a banned word in them.
Keep the attribution block that the source requires.

A contributor report and a pull request are not documentation.

## 4. Database Rules

Write a rule for each engine. Never write one rule for all engines.

### 4.1 Rules For Every Engine

* Do not execute raw SQL string concatenations.
* Use parameterized queries to prevent SQL injection.
* Wrap a multi-cell update in one transaction of the target engine.
* Quote each identifier through the `Dialect` port. Never quote it by hand.
* Never accept a table or column name from a request before the registry
  validates it.
* Keep each statement in the package for its engine.

### 4.2 Rules For The PostgreSQL Engine

* Use `pgx/v5` for the connection. Add another driver for another engine.
* Quote each identifier with `pgx.Identifier`.
* Wrap a multi-cell update with `BEGIN` and `COMMIT`.
* Run the migrations under a PostgreSQL advisory lock.

## 5. Build And Toolchain Rules

* Set `CGO_ENABLED=0` in every build, test, and release command.
* Add no dependency that needs CGO. Never use `go-sqlite3` or a C resolver.
* Keep `go.mod` free of any `require` block that pulls a C compiler.
* Run the engine locally with `air`. Read `air.toml` before you change it.
* Keep the local stack in `docker-compose.yml`.
* Keep no copy of the ASD-STE100 word lists. The vendored skill holds the
  Simplified Technical English rules.

## 6. Versioning And Release Rules

### 6.1 The Version Number

Kibtab uses semantic versioning. The number tells a user what changed.

| Part | Meaning |
| --- | --- |
| MAJOR | The public interface changed. A client must change. |
| MINOR | Kibtab added a feature. An old client keeps working. |
| PATCH | Kibtab fixed a defect. No interface changed. |

The public interface is the REST paths, the JSON field names, and the CLI
flags. A change to any of them needs a MAJOR version.

### 6.2 The Gate For A Release

* `CGO_ENABLED=0 go build ./...` succeeds with no environment variable.
* `CGO_ENABLED=0 go vet ./...` reports no problem.
* `CGO_ENABLED=0 go test ./...` passes.
* `make docs-check` and `make notice-check` report no problem.
* The module has no dependency that needs a C compiler.
* The code follows section 1. The docs follow section 3.
* The commits follow section 8. Each change follows section 9.

### 6.3 The Gate For A Ship

* Every check in section 6.2 passes.
* `make release-check` passes.
* The archive builds for each target with `CGO_ENABLED=0`.
* `ldd` reports no dynamic library for the Linux binary.
* `docs/changelogs/vX.Y.Z.md` exists. It states the upgrade steps.
* `docs/index.md` lists the release.
* The tag exists. The tag does not move.

### 6.4 The Release Files

* Release with GoReleaser. Read `.goreleaser.yaml` before you change it.
* Write one file at `docs/changelogs/vX.Y.Z.md` for each release.
* Copy `docs/changelogs/TEMPLATE.md`. Follow `docs/changelogs/README.md`.
* Update `docs/index.md` and `docs/CREDITS.md` with the release notes.
* Update `docs/THIRD_PARTY_NOTICES.md` when a dependency changes.
* Read `plan.md` for the scope of each version.

## 7. Repository Rules

* Commit no secret, no `.env` file, and no credential.
* Commit no build output, no `dist/` folder, and no coverage file.
* Keep the module path as `github.com/kibtab/kibtab`.
* Follow `plan.md` for scope. Never build a feature that a later version
  holds.
* Every task in `plan.md` uses a check box. Mark it done when the code and
  the release note exist. Delete the line when the release ships.

## 8. Commit Message Rules

Kibtab uses Conventional Commits. The format is
`<type>(<scope>): <description>`.
[CONTRIBUTING.md](CONTRIBUTING.md) holds the full table of types and scopes
with examples. Read it before your first commit.

* Keep the subject under 50 characters. Start with a verb.
* Do not end the subject with a full stop. Use British English.
* Choose the scope of the file that changed most. Never invent a scope.
* Use no scope when the change touches the whole repository.
* A change to a REST path, a JSON field name, or a CLI flag needs `!` after
  the scope and a `BREAKING CHANGE:` footer.
* Add a body when the subject does not explain the change. Wrap it at 72
  characters.
* Explain why the change was needed. Do not repeat what the diff shows.
* Link the issue at the end with `Closes #123`.
* Tag the commit. Never move a published tag.

After you finish a change, suggest a commit message.
Do not run the commit yourself unless the user asks for it.
Give the message in a fenced block. State why the scope fits the change.
Name the version in `plan.md` that the change belongs to.
Say whether the change needs a MAJOR version.

## 9. Scope And Separation Rules

Do one thing and do it well.
Give each file, each package, each commit, and each guide one job.

### 9.1 The Scope Of A Change

* Do one thing in a commit. Split the rest into another commit.
* Do one thing in a function. Split the rest into another function.
* Solve the problem in the task. Do not solve the next task.
* Never mix a repair with a feature, or a rename with a behaviour change.
* Never give a file, a package, or a function a second reason to change.
* Move code before you change it. Use one commit for the move.
* Add no flag, switch, or mode for a case that no user has.
* Add no abstraction before a second caller needs it.
* Keep the diff small. A large diff hides a large defect.
* Fail loudly and early. Never hide an error behind a default value.
* Prefer text a person can read over text only a tool can parse.

### 9.2 The Layer Rules

* Keep the domain free of transport, persistence, and presentation.
* Keep the transport free of SQL. Keep SQL free of the domain.
* Keep a configuration value out of the domain. Read it at the edge.
* Push the detail down to the layer that owns it.
* Read a stream and write a stream where the tool makes sense.
* Keep the core free of a convenience for any one adapter.

### 9.3 The Documentation Rules

* Give each guide one task. Copy no step into a second file.
* Put each setup step in the one file that owns it.
* Put a section of many files in its own folder. Add a `README.md` index.
* Link each file from the index. Leave no file unlinked.
* Update the index in the same commit as the file that it lists.
* Point to the document that holds the answer. Do not repeat it in a reply.
* Write the answer in one file. Link to it from every other place.
* Update a document in the same commit as the code it describes.

### 9.4 The Test Rules

* Test one concern per test. Split a test that asserts two things.
* Name each test for the concern that it covers.
* Keep the core test free of an adapter. Use a fake for a port.
* Run the contract suite for an adapter. Do not repeat it in each adapter.
* Keep a fixture in one file. Share it.
* Cover every line and every branch of the code you add or change.
* Aim for 100% coverage. Never lower the floor to make a build pass.
* Cover the error path. A happy-path test alone is not a test.
* Cover the boundary. Test the empty case and the overflow case.
* Use a table-driven test for each set of related cases.

### 9.5 The Coverage Rules

The module must reach 100% coverage.

* Run `make cover-verify`. It fails below the floor in `Makefile`.
* Run `make cover` for the per-function breakdown of each file.
* Write a test for each bug before you write the fix.
* Cover the code that a port supplies. A fake covers the port, not the
  adapter.
* Keep the core test suite free of a container. The engine suite uses a real
  database container.
* Never use `//go:build ignore` or a blank identifier to hide a line.
* Never skip a test to make a build pass. Fix the code or the test.

`COVERAGE_FLOOR` lives in `Makefile`.
Read [the architecture guide](docs/architecture.md) for the test layers.

### 9.6 The RTFM Rules

* Read the file before you change it. Read the whole file.
* Read `plan.md` before you start. Read this document before you edit.
* Search for the answer before you ask. Then search for the question.
* When a rule has no answer in the repository, write the rule first.

### 9.7 The Scope Test

Run this before you send a change for review.

1. Name the one thing the change does. Stop when you need a second verb.
2. Name the layer that owns the change. Stop when you need a second layer.
3. Name the version in `plan.md` that holds the change.
4. Read the diff. Remove each hunk that serves a second purpose.
5. Check that the commit message names that one thing.