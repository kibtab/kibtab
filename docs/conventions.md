# Code Conventions

Version 0.0.0-docs. This document holds the conventions for the Kibtab codebase.

Read [AGENTS.md](../AGENTS.md) for the rules that bind an agent.
Read [CONTRIBUTING.md](../CONTRIBUTING.md) for the steps to send a change.
Read [architecture.md](architecture.md) for the tree and the ports.

## Writing Code

* Use Go 1.22. Kibtab targets `go 1.22` in `go.mod`. Kibtab does not raise
  the version for a newer release. Go 1.22 is mature and it runs on every
  CI platform, on `golang:1.22-alpine`, and in the GoReleaser pipeline.
* Run `gofmt` on all code. Run `go vet` before you commit.
* Handle all errors directly. Never ignore a returned error.
* Use the standard library first. Add a dependency only for a real need.
* Keep functions under 50 lines. Keep files under 300 lines.
* Write table-driven tests for each set of related cases.
* Return a value. Never set a field for a caller in another layer.
* Fail loudly and early. Never hide an error behind a default value.
* Prefer text a person can read over text only a tool can parse.

## The Layer Rules

* Keep the domain free of transport, persistence, and presentation.
* Keep the transport free of SQL. Keep SQL free of the domain.
* Keep a configuration value out of the domain. Read it at the edge.
* Push the detail down to the layer that owns it.
* Read a stream and write a stream where the tool makes sense.
* Keep the core free of a convenience for any one adapter.

## Database Rules

Write a rule for each engine. Never write one rule for all engines.

### Rules For Every Engine

* Never execute raw SQL string concatenations.
* Use parameterized queries to prevent SQL injection.
* Wrap a multi-cell update in one transaction of the target engine.
* Quote each identifier through the `Dialect` port. Never quote it by hand.
* Never accept a table or column name from a request before the registry
  validates it.
* Keep each statement in the package for its engine.

### Rules For The PostgreSQL Engine

* Use `pgx/v5` for the connection. Add another driver for another engine.
* Quote each identifier with `pgx.Identifier`.
* Wrap a multi-cell update with `BEGIN` and `COMMIT`.
* Run the migrations under a PostgreSQL advisory lock.

## The Documentation Rules

Write every document in two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
Neither style replaces the other.
Both styles apply to every document, comment, and commit message.

Load the vendored skill first. Read
`../skills/simple-english/SKILL.md`.

### Simplified Technical English

Read [the catalogue](ste100/index.md) before you use a technical term.
It holds the approved Technical Names for Kibtab.

* Keep sentences shorter than 20 words.
* Use active voice and simple tenses.
* Use approved words from `../skills/simple-english/references/word-swaps.md`.
* Use one term for one meaning in each document.
* Put the condition before the command.
* Define each technical term at its first use.

### British English

* Use `-ise`, not `-ize`: `organise`, `recognise`, `analyse`.
* Use `-our`, not `-or`: `colour`, `behaviour`, `favour`.
* Use `-re`, not `-er`: `centre`, `metre`, `litre`.
* Use `licence` as a noun. Use `license` as a verb.
* Keep the proper name of a licence: `Apache License 2.0`.
* Double the final `l` after a vowel: `travelling`, `cancelled`.
* Use `artefact`, not `artifact`.
* Use `whilst` and `amongst` only where the sense needs them.
* Keep the plural `-ae` and `-a` forms: `formulae`.

### Rules For Both Styles

* Never use em dashes, contractions, or the words `should`, `would`, `may`,
  and `might`.
* Keep articles. Keep `that` after a verb.
* Put the version in the first line of every document.
* Run `make docs-check` before you commit.

## The Documentation Layout

* Give each guide one task. Copy no step into a second file.
* Put each setup step in the one file that owns it.
* Put a section of many files in its own folder. Add a `README.md` index.
* Link each file from the index. Leave no file unlinked.
* Update the index in the same commit as the file that it lists.
* Point to the document that holds the answer. Never repeat it in a reply.
* Write the answer in one file. Link to it from every other place.
* Update a document in the same commit as the code it describes.

[index.md](index.md) holds the table that names the owner of each step.

## The Canonical Documents

Three documents have one valid form. Copy the text as written.
[licence.md](licence.md) lists them and their canonical source.

Never apply the documentation styles to them.
Never fix a sentence, a long clause, or a banned word in them.
Keep each attribution block.

A contributor report and a pull request are not documentation.

## The Scope Of A Change

Do one thing and do it well.
Give each file, each package, each commit, and each guide one job.

* Do one thing in a commit. Split the rest into another commit.
* Do one thing in a function. Split the rest into another function.
* Solve the problem in the task. Never solve the next task.
* Never mix a repair with a feature, or a rename with a behaviour change.
* Never give a file, a package, or a function a second reason to change.
* Move code before you change it. Use one commit for the move.
* Add no flag, switch, or mode for a case that no user has.
* Add no abstraction before a second caller needs it.
* Keep the diff small. A large diff hides a large defect.

## The Scope Test

Run this before you send a change for review.

1. Name the one thing the change does. Stop when you need a second verb.
2. Name the layer that owns the change. Stop when you need a second layer.
3. Name the version in [plan.md](../plan.md) that holds the change.
4. Read the diff. Remove each hunk that serves a second purpose.
5. Check that the commit message names that one thing.

## The RTFM Rules

* Read the file before you change it. Read the whole file.
* Read [plan.md](../plan.md) before you start.
* Search for the answer before you ask. Then search for the question.
* When a rule has no answer in the repository, write the rule first.

## Build And Toolchain Rules

* Set `CGO_ENABLED=0` in every build, test, and release command.
* Add no dependency that needs CGO. Never use `go-sqlite3` or a C resolver.
* Keep `go.mod` free of any `require` block that pulls a C compiler.
* Run the engine locally with `air`. Read `air.toml` before you change it.
* Keep the local stack in `docker-compose.yml`.
* Keep no copy of the ASD-STE100 word lists. The vendored skill holds the
  Simplified Technical English rules.
* Keep the catalogue at [ste100/index.md](ste100/index.md). Add each approved
  name there before an author uses it.

## Next Steps

* Read [architecture.md](architecture.md) for the tree and the ports.
* Read [testing.md](testing.md) for the tests and the coverage.
* Read [releasing.md](releasing.md) for the version and the gate.
* Read [CONTRIBUTING.md](../CONTRIBUTING.md) to send a change.