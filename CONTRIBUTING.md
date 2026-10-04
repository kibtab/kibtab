# Contributing to Kibtab

Thank you for your interest in Kibtab.
Kibtab turns a spreadsheet into a client for a relational database.
Read this guide before you send a change.

By contributing, you agree to the
[Code of Conduct](CODE_OF_CONDUCT.md).

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

Read [AGENTS.md](AGENTS.md). It holds the rules for this repository.
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

Kibtab needs Go 1.22 or a newer version.
Kibtab builds without a C compiler.

```bash
CGO_ENABLED=0 go build ./...
air
```

The tool `air` reloads the engine on a code change.
Read `air.toml` before you change it.

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

Kibtab uses Conventional Commits.
Read [AGENTS.md](AGENTS.md) section 8 for the format and the scopes.

```text
feat(core): add the optimistic version compare
fix(postgres): stop the audit insert from rolling back the bump
docs(core): record the port rules in the architecture guide
```

## The Rules For A Change

* Do one thing in a commit. Split the rest into another commit.
* Do one thing in a function. Read [AGENTS.md](AGENTS.md) section 9.
* Keep the core free of an engine name and of a client name.
* Put each engine in its own package under `internal/adapters/driven/`.
* Put each transport in its own package under `internal/adapters/driver/`.
* Use a parameterized query. Never build a query by concatenation.
* Wrap a multi-cell update in one transaction.
* Set `CGO_ENABLED=0` in every build and test command.
* Write the documentation in Simplified Technical English and British English.
  Read [AGENTS.md](AGENTS.md) section 3 for both.

## Run The Checks

Run each check before you open a pull request.

```bash
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go vet ./...
CGO_ENABLED=0 go test ./...
make docs-check
make notice-check
```

Read the target list in [docs/index.md](docs/index.md).

## Pull Requests

* Open an issue first when the change is large or breaking.
* Keep each pull request to one purpose.
* Edit only the files that the change needs.
* Do not commit an editor metafile. Put it in your own global ignore file.
* Give the pull request a clear title and a short reason.
* Reference the issue with `Fixes #123`.
* Allow edits from the maintainers.
* Read your own diff before you send it.

## The Documentation

Each guide owns one task.
A step lives in one file only.
A section of many files has an index that links each file.

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