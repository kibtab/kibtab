# Kibtab

**Spreadsheet as client, database as server.**

Kibtab turns Microsoft Excel into a client for your relational database.
> And potentially other spreadsheet software in the future.

Your team keeps the familiar grid.
The engine keeps the transactions, the version checks, and the audit trail.

Kibtab rejects a write from an out-of-date sheet.
Every cell change is recorded.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
Read the [plan](plan.md) for the scope of each version.
No binary is published yet.

## Install

The install command works after the first release.

```bash
go install github.com/kibtab/kibtab/cmd/kibtab@latest
```

Kibtab builds without a C compiler. Every build sets `CGO_ENABLED=0`.

## Self-host

The stack holds Caddy, the engine, and PostgreSQL.

```bash
git clone https://github.com/kibtab/kibtab.git
cd kibtab
docker compose up -d
```

Read [the self-hosting guide](docs/self-hosting.md) after v0.8.0 ships.

## The Name

**Kibtab** joins two words.

*Kiban* (基盤) means *foundation* in Japanese.
It is the word for the base layer under a structure.

*Tabularium* is the Latin register of public tables.
It was the Roman archive where the state kept its records.

A foundation plus a register.
That is what this project is.

The name is pronounced `/kɪb-tæb/`.

## Licence

Apache License 2.0.
Read [LICENSE](LICENSE) for the full terms.

Kibtab ships third party software.
Read [the third party notices](docs/THIRD_PARTY_NOTICES.md) for each
dependency and its licence.

The engine licence finalises at v1.0.0. Read [the plan](plan.md) section 12.

## Documentation

| Document | Purpose |
| --- | --- |
| [Plan](plan.md) | The release plan. It splits the work by version. |
| [Documentation index](docs/index.md) | The list of every document. |
| [Agent rules](AGENTS.md) | The rules for AI coding agents. |
| [Changelog format](docs/changelogs/README.md) | The format for each release file. |
| [Credits](docs/CREDITS.md) | The people, the projects, and the sources. |
| [Architecture](docs/architecture.md) | The Hexagonal layers. Added at v0.2.0. |
| [REST interface](docs/openapi.yaml) | The API definition. Added at v0.5.0. |
| [Self-hosting](docs/self-hosting.md) | The install steps. Added at v0.8.0. |
| [Runbook](docs/runbook.md) | The operator steps. Added at v0.9.0. |

Documentation follows two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
Read [the agent rules](AGENTS.md) section 3 for both.

## Architecture

The code follows the Hexagonal architecture.
The core holds the domain logic. It has no database driver and no HTTP code.

```text
cmd/kibtab/                    The entry point and the CLI
internal/core/domain/          The models
internal/core/ports/           The interfaces
internal/core/services/        The business logic
internal/adapters/driven/      PostgreSQL with pgx/v5
internal/adapters/driver/      The HTTP handlers
```

Read [the architecture document](docs/architecture.md) at v0.2.0.

## Contributing

Read [AGENTS.md](AGENTS.md) before you write code.
Read [plan.md](plan.md) before you start a task.
Each task in the plan is a markdown check box.

Every commit uses Conventional Commits.
Read [AGENTS.md](AGENTS.md) section 8 for the format and the scopes.

## Acknowledgements

Built on the ideas in [docs/CREDITS.md](docs/CREDITS.md).
The documentation rules come from the vendored
[SimpleEnglish](https://github.com/AminBlg/SimpleEnglish) skill.
