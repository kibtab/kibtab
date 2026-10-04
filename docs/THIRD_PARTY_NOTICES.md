# Third Party Notices

Version 0.0.0-docs. Kibtab uses third party software.
This document lists that software and its licence.

Kibtab is distributed under the Apache License 2.0.
The Apache License 2.0 requires this document in each distribution.
GoReleaser adds this file to each archive.

## How To Update This Document

1. Read the licence of each new dependency.
2. Add a row to the table below.
3. Copy the licence text to `docs/licences/<module>@<version>.txt`.
4. Run `make notice-check` before you open a pull request.

The build fails when a dependency has no row in this file.

## Go Dependencies

The versions below are the plan for v0.1.0 to v1.0.0.
The list changes as the code changes.
The `go.mod` file holds the exact version.

| Module | Licence | Use in Kibtab | First version |
| --- | --- | --- | --- |
| `github.com/jackc/pgx/v5` | MIT | The PostgreSQL driver. | v0.3.0 |
| `github.com/go-chi/chi/v5` | MIT | The HTTP router. | v0.5.0 |
| `github.com/google/uuid` | BSD-3-Clause | The request and change identifiers. | v0.5.0 |
| `github.com/prometheus/client_golang` | Apache-2.0 | The metrics endpoint. | v0.9.0 |

The Go standard library carries a BSD-3-Clause licence.
Kibtab does not list it in this table.

## Client Dependencies

The Office.js client uses these packages.
They build in a separate CI job.

| Package | Licence | Use in Kibtab | First version |
| --- | --- | --- | --- |
| `office-js` | MIT | The Excel taskpane interface. | v0.6.0 |
| `typescript` | Apache-2.0 | The client build. | v0.6.0 |
| `esbuild` | MIT | The client bundle. | v0.6.0 |

## Build And Release Tools

These tools do not ship in the archive.
The licence still applies to a developer who installs them.

| Tool | Licence | Use in Kibtab | First version |
| --- | --- | --- | --- |
| `goreleaser` | MIT | The release build. | v0.8.0 |
| `air` | MIT | The hot reload for a local run. | v0.1.0 |
| `golangci-lint` | MIT | The Go lint. | v0.1.0 |
| `k6` | AGPL-3.0 | The load test. | v0.9.0 |

## Infrastructure Images

The local stack uses these images.
Each image carries its own licence.

| Image | Licence | Use in Kibtab | First version |
| --- | --- | --- | --- |
| `caddy:2-alpine` | Apache-2.0 | The edge proxy and the TLS. | v0.8.0 |
| `postgres:16-alpine` | PostgreSQL License | The database in the local stack. | v0.8.0 |

## Vendored Skills

| Project | Licence | Use in Kibtab | Pinned commit |
| --- | --- | --- | --- |
| [SimpleEnglish](https://github.com/AminBlg/SimpleEnglish) | MIT | The documentation rules. | `32ea2d3f4404bfbff162c1861e5ae4f1bcc628b3` |

The copy lives in `skills/simple-english/`.
Read `skills/simple-english/VENDOR.md` for the refresh command.
The licence text lives in `skills/simple-english/LICENSE`.

## Standard Notices

Apache License 2.0, section 4

> If the Work includes a NOTICE text file as part of its distribution, then any
> Derivative Works that You distribute must include a readable copy of the
> attribution notices contained within such NOTICE file.

Kibtab keeps this requirement for each dependency with an Apache-2.0 licence.