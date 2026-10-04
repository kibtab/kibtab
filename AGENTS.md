# Agent Rules for Kibtab

This document gives instructions for AI coding agents.
You must follow these rules when you write code or edit files in this repository.

## 1. Code Architecture Rules

* Keep code in the correct Hexagonal layer.
* Do not import database drivers into `internal/core/`.
* `internal/core/` must contain only pure Go domain logic and ports.
* Put all PostgreSQL queries inside `internal/adapters/driven/postgres/`.
* Put all HTTP handlers inside `internal/adapters/driver/http/`.
* Put all CLI commands inside `cmd/kibtab/`.
* Define interfaces in `internal/core/ports/`. Implement them in adapters.

## 2. Writing Code

* Use Go 1.22 or a newer version.
* Write clear Go code. Use `gofmt` to format all code.
* Handle all errors directly. Do not ignore returned errors.
* Use `pgx/v5` for PostgreSQL connections.
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
* Run `make docs-check` before you commit. The target checks the sentence
  length and the banned words. It does not check the word lists.

## 4. Database Rules

* Do not execute raw SQL string concatenations.
* Use parameterized queries to prevent SQL injection.
* Always wrap multi-cell updates in a PostgreSQL transaction (`BEGIN` and
  `COMMIT`).
* Quote identifiers with `pgx.Identifier` when a query needs a table name.
* Never accept a table or column name from a request without validating it
  against the table registry.

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
| `postgres` | `internal/adapters/driven/postgres/`. |
| `http` | `internal/adapters/driver/http/`. The routes and the codec. |
| `cli` | `cmd/kibtab/`. The entry point and the flags. |
| `client` | `client/`. The Office.js taskpane. |
| `schema` | The `_kibtab_meta` schema and the migrations. |
| `deploy` | The `Dockerfile`, the `Caddyfile`, and the compose file. |
| `release` | GoReleaser, the tags, and the archives. |
| `docs` | Every file under `docs/`, plus `README.md` and `AGENTS.md`. |
| `deps` | `go.mod`, `go.sum`, and the client packages. |
| `ci` | The workflow files. |
| `skills` | The vendored skill in `skills/simple-english/`. |

### 8.4 Examples

```text
feat(core): add the optimistic version compare
fix(postgres): stop the audit insert from rolling back the version bump
docs(core): record the port rules in the architecture guide
feat(http): add the POST /v1/tables/{table}/cells route
refactor(postgres): move the row query out of the table registry
build(release): set CGO_ENABLED=0 in the goreleaser environment
test(postgres): cover the rollback of a failed batch
chore(skills): pin the simple-english skill to 32ea2d3
```

### 8.5 Scope Selection Rules

* Choose the scope of the file that changed most.
* Use `deps` for a version bump. Use `release` for a build setting.
* Use `core` for a change in more than one package under `internal/core/`.
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