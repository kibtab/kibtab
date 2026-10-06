<!--
Badges. Uncomment these lines once `make coverage-svg` writes coverage.svg.
The coverage badge needs a test run first. Keep the badge links in step
with the repository name.

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kibtab/kibtab?style=for-the-badge&logo=go)
![GitHub License](https://img.shields.io/github/license/kibtab/kibtab?style=for-the-badge)

![Coverage](coverage.svg)
-->

# Kibtab

**Spreadsheet as client, database as server.**

Kibtab turns a spreadsheet into a client for your relational database.
> Excel is the first client. Other spreadsheet software can come later.

Your team keeps the familiar grid.
Kibtab keeps the transactions, the version checks, and the audit trail.

Kibtab rejects a write from an out-of-date sheet.
Every cell change is recorded.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
No binary is published yet.
Read [the plan](https://kibtab.readthedocs.io/en/latest/#releases) for the scope of each version.
Or read [the plan file](plan.md) in the repository.

## The Name

**Kibtab** joins two words.

*Kiban* (基盤) means *foundation* in Japanese.
It is the word for the base layer under a structure.

*Tabularium* is the Latin register of public tables.
It was the Roman archive where the state kept its records.

A foundation plus a register.
That is what this project is.

The name is pronounced `/kɪb-tæb/`.

## Getting Started

Read the guide for your task.

| Task | Guide |
| --- | --- |
| Install a binary or build from source | [the install guide](https://kibtab.readthedocs.io/en/latest/install.html) |
| Run the stack on a server | [the self-hosting guides](https://kibtab.readthedocs.io/en/latest/self-hosting/README.html) |
| Read the layers and the port contract | [the architecture guide](https://kibtab.readthedocs.io/en/latest/architecture.html) |
| Read the terms | [the licence guide](https://kibtab.readthedocs.io/en/latest/licence.html) |

## Architecture

The code follows the Hexagonal architecture.
The core holds the domain logic. It has no database driver and no HTTP code.
The core names no database engine and no spreadsheet client.

```text
cmd/kibtab/                    The wiring and the CLI
internal/core/domain/          The models. They name no engine or client.
internal/core/ports/           The interfaces. The core owns them.
internal/core/services/        The use cases. They call the ports only.
internal/adapters/driven/      One package per database engine.
internal/adapters/driver/      One package per transport.
client/                        One folder per spreadsheet client.
```

A new database means one new package under `adapters/driven/`.
A new spreadsheet means one new client folder.
Neither changes a file under `internal/core/`.

Read [the architecture guide](https://kibtab.readthedocs.io/en/latest/architecture.html) for the full design.

## Documentation

| Document | Purpose |
| --- | --- |
| [Plan](https://kibtab.readthedocs.io/en/latest/#releases) | The release plan. It splits the work by version. |
| [Documentation index](https://kibtab.readthedocs.io/en/latest/) | The list of every document. |
| [Agent rules](AGENTS.md) | The rules for AI coding agents. |
| [Conventions](https://kibtab.readthedocs.io/en/latest/conventions.html) | The code, the documentation, and the scope rules. |
| [Technical names](https://kibtab.readthedocs.io/en/latest/ste100/index.html) | The approved technical names for Kibtab. |
| [Testing](https://kibtab.readthedocs.io/en/latest/testing.html) | The kinds of test and the coverage rules. |
| [Releasing](https://kibtab.readthedocs.io/en/latest/releasing.html) | The version number, the gate, and the release steps. |
| [Licence](https://kibtab.readthedocs.io/en/latest/licence.html) | The terms for Kibtab and the skill. |
| [Third party notices](https://kibtab.readthedocs.io/en/latest/THIRD_PARTY_NOTICES.html) | The licences of each dependency. |
| [Changelog format](https://kibtab.readthedocs.io/en/latest/changelogs/README.html) | The format for each release file. |
| [Credits](https://kibtab.readthedocs.io/en/latest/CREDITS.html) | The people, the projects, and the sources. |
| [Contributing](CONTRIBUTING.md) | The steps to send a change. |
| [Code of Conduct](CODE_OF_CONDUCT.md) | The terms for taking part. |
| [Security](SECURITY.md) | The private route for a vulnerability. |

Documentation follows two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
Read section 3 of [the agent rules](AGENTS.md) for both.

## Contributing

Read [the agent rules](AGENTS.md) before you write code.
Read [the plan](https://kibtab.readthedocs.io/en/latest/#releases) before you start a task.
Each task in the plan is a markdown check box.
Read [the contributing guide](CONTRIBUTING.md) for the full steps.

Every commit uses Conventional Commits.
Read section 6 of [the agent rules](AGENTS.md) for the format and the scopes.

## Acknowledgements

Built on the ideas in [the credits](https://kibtab.readthedocs.io/en/latest/CREDITS.html).
The documentation rules come from the vendored
[SimpleEnglish](https://github.com/AminBlg/SimpleEnglish) skill.
