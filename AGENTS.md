# Agent Context for Kibtab

This document holds the context that an AI coding agent needs.
The rules for a change live in [the contributing guide](CONTRIBUTING.md).
A person and an agent follow the same rules, so a rule is written once.
A person is the primary reader of the contributing guide.

This file adds the agent context and the routing. It repeats no rule.

## Read First

| Read this | For |
| --- | --- |
| [the contributing guide](CONTRIBUTING.md) | The commit format, the scopes, the DCO, the tests, and the pull request. |
| [the plan](plan.md) | The scope of each version. The open check boxes. |
| [the conventions guide](docs/conventions.md) | The code, the documentation, and the scope rules. |
| [the catalogue](docs/ste100/index.md) | The approved Technical Names. Use them in every document. |
| [the architecture guide](docs/architecture.md) | The tree, the ports, and the data flow. |
| [the testing guide](docs/testing.md) | The unit tests, the end-to-end tests, and the coverage. |
| [the releasing guide](docs/releasing.md) | The version number, the gate, and the release steps. |
| [the licence guide](docs/licence.md) | The licence terms and the canonical documents. |
| [the install guide](docs/install.md) | Installing and building. |
| [the self-hosting guides](docs/self-hosting/README.md) | Running Kibtab on a server. |

Read the linked document before you act. Do not work from this file alone.

## The Working Agreement

* Do the work for one commit only. Stop when that work is complete.
* Read `plan.md` before you start. Pick one open box.
* Mark the box in the same commit as the work that it records.
* Write the release note in `docs/changelogs/` in the same commit.
* Name the version in `plan.md` and the MAJOR answer in the reply.
* Give the commit message in a fenced block. Do not run the commit.
* Wait for the user before you start the next change.

## Where Each Rule Lives

A rule has one home. Read that home before you act.

| Rule | Document |
| --- | --- |
| The commit format, the scope table, and the DCO | [the contributing guide](CONTRIBUTING.md) |
| The code rules and the scope test | [the conventions guide](docs/conventions.md) |
| The writing styles and the documentation layout | [the conventions guide](docs/conventions.md) |
| The approved Technical Names | [the catalogue](docs/ste100/index.md) |
| The layers, the ports, and the tree | [the architecture guide](docs/architecture.md) |
| The tests and the coverage floor | [the testing guide](docs/testing.md) |
| The version number and the release gate | [the releasing guide](docs/releasing.md) |

## Reading The Repository

* Read the file before you change it. Read the whole file.
* Search for the answer before you ask. Then search for the question.
* When a rule has no answer in the repository, write the rule first.
* Write the answer in one file. Link to it from every other place.
