# Credits

Version 0.0.0-docs. This document lists the people and projects behind Kibtab.
It also lists the projects that Kibtab learns from.

Kibtab is an open core project.
Every credit below is voluntary.
Add your name in the pull request that adds your work.

## Maintainers

| Name | Role |
| --- | --- |
| Kibtab contributors | Kibtab, the clients, and the documentation. |

Find the current list in the `AUTHORS` file and in the repository contributors.

## Inspiration

Kibtab does not copy the projects below.
These projects show what Kibtab builds on.

| Project | Why it matters |
| --- | --- |
| [Hexagonal architecture](https://alistair.cockburn.us/hexagonal-architecture/) | The layer rules in `AGENTS.md`. |
| [PostgreSQL](https://www.postgresql.org/) | The primary database. It holds the transactions. |
| [pgx](https://github.com/jackc/pgx) | The Go driver. It supports a parameterized query. |
| [Microsoft Office Add-ins](https://learn.microsoft.com/office/dev/add-ins/) | The taskpane platform for Excel. |
| [Google Apps Script](https://developers.google.com/apps-script) | The platform for the future Sheets adapter. |
| [Caddy](https://caddyserver.com/) | The edge proxy. It manages the TLS. |
| [GoReleaser](https://goreleaser.com/) | The release build for each target. |
| [Air](https://github.com/air-verse/air) | The hot reload for a local run. |

## Standards

| Standard | Why it matters |
| --- | --- |
| [ASD-STE100 Issue 9](https://www.asd-ste100.org/) | The Simplified Technical English rules for the documentation. |
| [Semantic Versioning](https://semver.org/) | The version number in `plan.md`. |

## Documentation Rules

The `skills/simple-english/` directory holds a copy of the
[SimpleEnglish](https://github.com/AminBlg/SimpleEnglish) skill by
[AminBlg](https://github.com/AminBlg). The skill carries an MIT licence.

The vendored skill gives Kibtab the word rules in `../AGENTS.md` section 3.
Read `skills/simple-english/VENDOR.md` for the pinned commit.

## Icons And Assets

| Asset | Source | Licence |
| --- | --- | --- |
| The Kibtab logo | Made for this project. | Apache-2.0. |
| The README diagram | Made with Mermaid. | Mermaid is MIT. |

## Test Data

Kibtab tests use these sample data sets.

| Sample | Source | Licence |
| --- | --- | --- |
| The sample sales table | Made for this project. | Apache-2.0. |
| The conflict test fixture | Made for this project. | Apache-2.0. |

No test data holds personal data.

## How To Add A Credit

1. Add a row to the table for the section.
2. Use the name that the person asks for.
3. Keep the table in alphabetical order by name.
4. Update this file in the same commit as the work.

Record the work in `docs/changelogs/vX.Y.Z.md` as well.