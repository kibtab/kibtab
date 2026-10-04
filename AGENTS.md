# Agent Rules for Kibtab

This document gives instructions for AI coding agents.
You must follow these rules when you write code or edit files in this repository.

`plan.md` holds the work. `docs/` holds the detail.
This document holds the rules that bind you. It points at the detail.

| Read this | For |
| --- | --- |
| [the plan](plan.md) | The scope of each version. The open check boxes. |
| [the conventions guide](docs/conventions.md) | The code, docs, and scope conventions. |
| [the catalogue](docs/ste100/index.md) | The approved Technical Names. Use them in every document. |
| [the architecture guide](docs/architecture.md) | The tree, the ports, and the data flow. |
| [the testing guide](docs/testing.md) | The unit tests and the coverage. |
| [the releasing guide](docs/releasing.md) | The version, the gate, and the release. |
| [the contributing guide](CONTRIBUTING.md) | The commit format, the scope table, and the DCO. |
| [the licence guide](docs/licence.md) | The licence terms and the canonical documents. |
| [the install guide](docs/install.md) | Installing and building. |
| [the self-hosting guides](docs/self-hosting/README.md) | Running Kibtab on a server. |

Read the linked document before you act. Do not work from this file alone.

## 1. Architecture

`internal/core/` holds pure Go domain logic and ports. It has no I/O.
It names no database engine and no spreadsheet client.

* Define each interface in `internal/core/ports/`. Implement it in an adapter.
* Put each engine in its own package under `internal/adapters/driven/`.
* Put each transport in its own package under `internal/adapters/driver/`.
* Put each spreadsheet client in its own folder under `client/`.
* Build each adapter in one place, under `cmd/kibtab/`.
* Reach a database only through a port.
* Use the `Clock` port in the core. Never call the system clock in a test.
* Never put a client name, such as Excel or Sheets, in `internal/core/`.
* Add a CI step that fails the build when the core imports an adapter.

Read [the architecture guide](docs/architecture.md) for the tree, the port
table, the data flow, and the steps to add an adapter.

## 2. Code

* Use Go 1.22. Set `go 1.22` in `go.mod`. Run `gofmt` and `go vet`.
* Handle all errors directly. Never ignore a returned error.
* Set `CGO_ENABLED=0` in every build, test, and release command.
* Add no dependency that needs CGO.
* Never execute raw SQL string concatenations. Use a parameterized query.
* Wrap a multi-cell update in one transaction.
* Quote each identifier through the `Dialect` port.
* Never accept a table or column name from a request before the registry
  validates it.
* Keep `internal/core/` free of build tags. A build tag changes an adapter,
  never the domain.

Read [the conventions guide](docs/conventions.md) for the full rules.

## 3. Documentation

Write every document in two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
Both styles apply to every document, comment, and commit message.

* Keep sentences shorter than 20 words. Use active voice.
* Use each Technical Name as the catalogue defines it.
  Add a name to the catalogue before you use it.
* Use `-ise` and `-our`, never `-ize` or `-or`. Use `licence` as a noun.
* Never use em dashes, contractions, or the words `should`, `would`, `may`,
  and `might`.
* Put the version in the first line of every guide under `docs/`.
  Put no version line at the repository root.
* Run `make docs-check` before you commit.

Three documents have one valid form. Copy them as written.
[the licence guide](docs/licence.md) lists them.

Read [the conventions guide](docs/conventions.md) for the word rules and the
documentation layout.
Read [the catalogue](docs/ste100/index.md) before you use a technical
term. It holds the approved names and the rejected words.

## 4. Tests

The module must reach 100% coverage.

* Write a test for each bug before you write the fix.
* Cover the error path and the boundary. Cover every branch you touch.
* Keep the core test free of an adapter. Use a fake for a port.
* Run the contract suite for an adapter. Never repeat it in each adapter.
* Never lower the floor to make a build pass. Never skip a test to pass.

```bash
make cover-verify
```

Read [the testing guide](docs/testing.md) for the test layers and the suites.

## 5. Versioning And Release

| Part | Meaning |
| --- | --- |
| MAJOR | The public interface changed. A client must change. |
| MINOR | Kibtab added a feature. An old client keeps working. |
| PATCH | Kibtab fixed a defect. No interface changed. |

The public interface is the REST paths, the JSON field names, and the CLI
flags. A change to any of them needs a MAJOR version and a `!` after the
scope.

* Release with GoReleaser. Read `.goreleaser.yaml` before you change it.
* Write one file at `docs/changelogs/vX.Y.Z.md` for each release.
* Tag the commit. Never move a published tag.

Read [the releasing guide](docs/releasing.md) for the gate and the steps.

## 6. Commits

Kibtab uses Conventional Commits. The format is
`<type>(<scope>): <description>`.

* Keep the subject under 50 characters. Start with a verb.
* Do not end the subject with a full stop. Use British English.
* Name the scope of the component that the change serves.
  Use `db-postgres`, `client-excel`, or `transport-http`.
  Never use a bare engine name or a bare client name.
* Use no scope when the change touches the whole repository.
* Split a change that touches two components into two commits.
* Never invent a scope. Add it to the table in `CONTRIBUTING.md` first.
* Explain why the change was needed. Never repeat what the diff shows.
* Link the issue at the end with `Closes #123`.

Sign off each commit. Read [the contributing guide](CONTRIBUTING.md) for the DCO.

After you finish a change, suggest a commit message.
Do not run the commit yourself unless the user asks for it.
Give the message in a fenced block. State why the scope fits the change.
Name the version in `plan.md` that the change belongs to.
Say whether the change needs a MAJOR version.

[the contributing guide](CONTRIBUTING.md) holds the type table and the scope table.
[the conventions guide](docs/conventions.md) holds the full scope rules.
Read the scope table before you write a commit.

## 7. Scope

Do one thing and do it well.
Give each file, each package, each commit, and each guide one job.

* Do one thing in a commit. Split the rest into another commit.
* Solve the problem in the task. Never solve the next task.
* Add no flag, switch, or mode for a case that no user has.
* Keep the diff small. A large diff hides a large defect.
* Follow `plan.md` for scope. Never build a feature that a later version
  holds.
* Every task in `plan.md` uses a check box. Mark it done when the code and
  the release note exist. Delete the line when the release ships.
* Commit no secret, no `.env` file, and no build output.
* Keep the module path as `github.com/kibtab/kibtab`.

Read [the conventions guide](docs/conventions.md) for the scope test and the
layer rules.

## 8. Reading The Repository

* Read the file before you change it. Read the whole file.
* Read `plan.md` before you start. Read the linked document before you edit.
* Search for the answer before you ask. Then search for the question.
* Point to the document that holds the answer. Never repeat it in a reply.
* Write the answer in one file. Link to it from every other place.
* When a rule has no answer in the repository, write the rule first.