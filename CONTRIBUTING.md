# Contributing to Kibtab

Thank you for your interest in Kibtab.
Kibtab turns a spreadsheet into a client for a relational database.
Read this guide before you send a change.

By contributing, you agree to the
[Code of Conduct](CODE_OF_CONDUCT.md).

## Where To Read

Read the document that matches your role.
Each document holds its own rules. Do not work from this guide alone.

| Reader | Read |
| --- | --- |
| A person | [The documentation index](docs/index.md). It lists every guide. |
| An AI agent | [AGENTS.md](AGENTS.md). It holds the rules that bind an agent. |
| A person who commits | [The commit format](#the-commit-message) below. |
| A person who tests | [docs/testing.md](docs/testing.md). |
| A person who releases | [docs/releasing.md](docs/releasing.md). |
| A person who writes words | [docs/conventions.md](docs/conventions.md). |
| A person who adds a term | [docs/ste100/index.md](docs/ste100/index.md). |

This guide holds the steps to send a change.
The [documentation index](docs/index.md) holds the detail of each topic.
The [agent rules](AGENTS.md) hold the rules that bind an AI agent.

## Ways To Contribute

Code is not the only way to help.
These contributions all matter.

* Fix a defect.
* Add a feature.
* Improve a document.
* Triage an issue.
* Review a pull request.
* Answer a question in an issue.

## Before You Start

Read [plan.md](plan.md). It holds the scope of each version.

Every task in the plan is a markdown check box.
Pick an open box. Read the version that owns it.
Stay inside that version.

Read [the architecture guide](docs/architecture.md) before you move code.
It holds the layer rules and the port contract.

## Set Up

Clone the repository and fetch the toolchain.

```bash
git clone https://github.com/kibtab/kibtab.git
cd kibtab
go version
```

Kibtab needs Go 1.22.
Kibtab builds without a C compiler.

```bash
make build
make watch
```

The target `make watch` reloads the engine on a code change.
Read `air.toml` before you change it.
Run `make help` for the full list of targets.

## Issues

Search the open and closed issues before you open a new one.

* Do not bump an issue that needs no answer.
* For a defect, give the steps, the version, and the output.
* For a defect, give a small case that shows the fault.
* Do not post a security issue in a public issue. Read [SECURITY.md](SECURITY.md).

## The Changelog

Each release has one file at `docs/changelogs/vX.Y.Z.md`.
Read [the changelog format](docs/changelogs/README.md) before you write one.
Copy [the template](docs/changelogs/TEMPLATE.md) first.

## The Commit Message

Kibtab uses Conventional Commits. The format is
`<type>(<scope>): <description>`.

| Type | Use for |
| --- | --- |
| `feat` | A new feature for a user. |
| `fix` | A defect repair. |
| `docs` | A documentation change only. |
| `refactor` | A code change with no new behaviour. |
| `test` | A test change only. |
| `perf` | A change that improves a measured result. |
| `build` | The toolchain, a dependency, or a build setting. |
| `ci` | The workflow files. |
| `chore` | A task with no source change. |

### The Scope Names A Component

Kibtab serves any database and any spreadsheet.
A commit that names only a folder hides which side of the product it changes.
The scope names the component instead.
A component scope has two parts. It names the role and then the name.

| Part | Meaning |
| --- | --- |
| `db-<engine>` | One database engine. For example `db-postgres`. |
| `client-<spreadsheet>` | One spreadsheet client. For example `client-excel`. |
| `transport-<name>` | One transport. For example `transport-http`. |

Use a two-part scope for every engine, every client, and every transport.
Add the name to the table below when the component first lands.
Use the name of the folder in the repository.

### The Scope Table

| Scope | Covers |
| --- | --- |
| `core` | `internal/core/`. The domain, the ports, the services. |
| `db-postgres` | `internal/adapters/driven/postgres/`. |
| `db-duckdb` | `internal/adapters/driven/duckdb/`. |
| `transport-http` | `internal/adapters/driver/http/`. |
| `client-excel` | `client/excel/`. |
| `client-sheets` | `client/sheets/`. |
| `client-local` | `client/local/`. |
| `schema` | The `_kibtab_meta` schema and the migrations. |
| `cli` | `cmd/kibtab/`. The entry point and the flags. |
| `wiring` | The adapter wiring in `cmd/kibtab/`. |
| `deploy` | The `Dockerfile`, the `Caddyfile`, and the compose file. |
| `release` | GoReleaser, the tags, and the archives. |
| `build` | The `Makefile`, `air.toml`, and `.gitignore`. |
| `docs` | Every file under `docs/`, plus `README.md` and `AGENTS.md`. |
| `deps` | `go.mod` and `go.sum`. |
| `ci` | The workflow files. |
| `skills` | The vendored skill in `skills/simple-english/`. |

### The Rules For A Scope

* Keep the subject under 50 characters. Start with a verb.
* Do not end the subject with a full stop.
* Name the one component that the change serves.
* Use no scope when the change touches the whole repository.
* Never write a bare engine name such as `postgres` as the scope.
  Write `db-postgres`.
* Never write a bare client name such as `excel` as the scope.
  Write `client-excel`.
* Never invent a scope. Add it to the table first.
* Split a change that touches two components into two commits.

A core change stays `core`, whatever adapters read it.
A change in two adapters is two commits.
One commit never carries a bare name that hides the role.

```text
feat(core): add the optimistic version compare
fix(db-postgres): stop the audit insert from rolling back the bump
feat(transport-http): add the POST cells route
feat(db-duckdb): add the DuckDB dialect behind the Dialect port
feat(client-excel): map the frozen range onto the sheet
build(release): set CGO_ENABLED=0 in the goreleaser environment
test(db-postgres): cover the rollback of a failed batch
chore(skills): pin the simple-english skill to 32ea2d3
```

A change to a REST path, a JSON field name, or a CLI flag is breaking.
Write `!` after the scope. Then add a `BREAKING CHANGE:` footer.

```text
feat(transport-http)!: version the cell write body

The v1 field names move to the envelope meta object.
A client must read the version from meta.version.

BREAKING CHANGE: the v1 JSON field names are removed.
Read docs/changelogs/v2.0.0.md before you upgrade.
```

## The Tests

The module aims for 100% coverage. Every line and every branch needs a test.

Read [the testing guide](docs/testing.md) for the full rules.
It holds the kinds of test, the suites, and the coverage rules.

Read the coverage before you send the change.

```bash
make cover
make cover-verify
```

`make cover-verify` fails below the floor in `Makefile`.
The floor is 100.
Never lower the floor to make a build pass.
Never skip a test to make a build pass.

Read [AGENTS.md](AGENTS.md) section 4 for the rules that bind an agent.

## The Rules For A Change

Read [the conventions guide](docs/conventions.md) for the full rules.
It holds the code rules, the layer rules, and the database rules.

The rules that most often catch a change:

* Keep `internal/core/` free of an engine name and of a client name.
* Put each engine in its own package under `internal/adapters/driven/`.
* Put each transport in its own package under `internal/adapters/driver/`.
* Use a parameterized query. Never build a query by concatenation.
* Wrap a multi-cell update in one transaction.
* Set `CGO_ENABLED=0` in every build and test command.
* Write the documentation in Simplified Technical English and British English.
  Read [AGENTS.md](AGENTS.md) section 3 for both.

Read [AGENTS.md](AGENTS.md) section 7 for the scope rules.

## Run The Checks

Run `make check` before you open a pull request.
It runs the build, the vet, the tests, and both documentation gates.

```bash
make check
```

Run each gate on its own when you need one.

```bash
make vet
make lint
make test
make cover-verify
make docs-check
make notice-check
```

Read `make help` for the full list of targets.
Read the gate rules in [AGENTS.md](AGENTS.md) section 5.

## Pull Requests

* Open an issue first when the change is large or breaking.
* Keep each pull request to one purpose.
* Edit only the files that the change needs.
* Do not commit an editor metafile. Put it in your own global ignore file.
* Give the pull request a clear title and a short reason.
* Reference the issue with `Fixes #123`.
* Allow edits from the maintainers.
* Read your own diff before you send it.

## Developer Certificate of Origin (DCO)

We encourage all contributors to sign off on their commits using the `-s` or
`--signoff` flag with `git commit`:

```bash
git commit -s -m "feat(db-postgres): add the batch delta validation"
```

Signing off indicates that you have the right to submit your contribution
under the Apache License 2.0. While we do not automatically reject pull
requests missing a sign-off, adopting this habit helps us maintain clear
provenance across the codebase as Kibtab grows.

Set the flag once for this repository to sign every commit.

```bash
git config --local commit.signoff true
```

The sign-off line sits under the commit message.

```text
feat(db-postgres): add the batch delta validation

Signed-off-by: Your Name <you@example.com>
```

## The Documentation

Each guide owns one task.
A step lives in one file only.
A section of many files has an index that links each file.

Read [the conventions guide](docs/conventions.md) for the documentation rules.
It holds the documentation layout and the writing styles.
Read [the documentation index](docs/index.md) for the list of each guide.

* Add a guide beside the task that it covers.
* Add the link to the index in the same commit.
* Update the guide in the same commit as the code that it describes.

## Review

* Push each commit. Do not squash. A reviewer reads each step.
* Present a solution. Do not send a question with no code.
* Read the diff after each new commit.
* Be patient. The maintainer reads each pull request.

## Security

Do not post a security issue in a public issue.
Read [SECURITY.md](SECURITY.md) for the private route.

## Licence

A contribution uses the Apache License 2.0.
Do not send a licence change.
Read [the licence guide](docs/licence.md) for the terms.