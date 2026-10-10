# Third Party Notices

Version 0.0.0-docs. Kibtab uses third party software.
This document lists that software and its licence.

Kibtab is distributed under the Apache License 2.0.
The Apache License 2.0 requires this document in each distribution.
GoReleaser adds this file to each archive.

## How To Update This Document

1. Read the licence of each new dependency.
2. Add a row to the table below.
3. Link the licence text in the upstream repository.
4. Run `make notice-check` before you open a pull request.

The build fails when a dependency has no row in this file.

## The Link Rule

Kibtab links to each licence text in the upstream repository.
Kibtab copies no third party licence into this repository.
The upstream copy stays the one source.

A copy goes stale when a dependency changes its licence.
A copy also hides a wrong version from a reader.

Pin each link to the tag in `go.mod`.
Use the repository branch for a dependency that no release holds yet.

The copy under `skills/simple-english/` stays.
It holds source code, not only a licence.

## Go Dependencies

The versions below are the plan for v0.1.0 to v1.0.0.
The list changes as the code changes.
The `go.mod` file holds the exact version.

| Module | Licence | Use in Kibtab | First version | Licence text |
| --- | --- | --- | --- | --- |
| `github.com/jackc/pgx/v5` | MIT | The PostgreSQL driver. | v0.3.0 | [LICENSE](https://github.com/jackc/pgx/blob/v5.6.0/LICENSE) |
| `github.com/jackc/puddle/v2` | MIT | The connection pool for pgx. | v0.3.0 | [LICENSE](https://github.com/jackc/puddle/blob/v2.2.1/LICENSE) |
| `github.com/jackc/pgpassfile` | MIT | The password file lookup for pgx. | v0.3.0 | [LICENSE](https://github.com/jackc/pgpassfile/blob/v1.0.0/LICENSE) |
| `github.com/jackc/pgservicefile` | MIT | The service file lookup for pgx. | v0.3.0 | [LICENSE](https://github.com/jackc/pgservicefile/blob/091c0ba34f0a/LICENSE) |
| `golang.org/x/crypto` | BSD-3-Clause | The SSH and password hashing for pgx. | v0.3.0 | [LICENSE](https://cs.opensource.google/go/x/crypto/+/refs/tags/v0.17.0:LICENSE) |
| `golang.org/x/sync` | BSD-3-Clause | The concurrency primitives for pgx. | v0.3.0 | [LICENSE](https://cs.opensource.google/go/x/sync/+/refs/tags/v0.1.0:LICENSE) |
| `golang.org/x/text` | BSD-3-Clause | The locale data for pgx. | v0.3.0 | [LICENSE](https://cs.opensource.google/go/x/text/+/refs/tags/v0.14.0:LICENSE) |
| `github.com/go-chi/chi/v5` | MIT | The HTTP router. | v0.5.0 | [LICENSE](https://github.com/go-chi/chi/blob/master/LICENSE) |
| `github.com/google/uuid` | BSD-3-Clause | The request and change identifiers. | v0.5.0 | [LICENSE](https://github.com/google/uuid/blob/master/LICENSE) |
| `github.com/prometheus/client_golang` | Apache-2.0 | The metrics endpoint. | v0.9.0 | [LICENSE](https://github.com/prometheus/client_golang/blob/main/LICENSE) |

The Go standard library carries a BSD-3-Clause licence.
Kibtab does not list it in this table.

## Client Dependencies

The Office.js client uses these packages.
They build in a separate CI job.

| Package | Licence | Use in Kibtab | First version | Licence text |
| --- | --- | --- | --- | --- |
| `office-js` | MIT | The Excel taskpane interface. | v0.6.0 | [LICENSE](https://github.com/OfficeDev/office-js/blob/release/LICENSE.md) |
| `typescript` | Apache-2.0 | The client build. | v0.6.0 | [LICENSE](https://github.com/microsoft/TypeScript/blob/main/LICENSE.txt) |
| `esbuild` | MIT | The client bundle. | v0.6.0 | [LICENSE](https://github.com/evanw/esbuild/blob/main/LICENSE.md) |

## Build And Release Tools

These tools do not ship in the archive.
The licence still applies to a developer who installs them.

| Tool | Licence | Use in Kibtab | First version | Licence text |
| --- | --- | --- | --- | --- |
| `goreleaser` | MIT | The release build. | v0.8.0 | [LICENSE](https://github.com/goreleaser/goreleaser/blob/main/LICENSE.md) |
| `air` | MIT | The hot reload for a local run. | v0.1.0 | [LICENSE](https://github.com/air-verse/air/blob/master/LICENSE) |
| `golangci-lint` | MIT | The Go lint. | v0.1.0 | [LICENSE](https://github.com/golangci/golangci-lint/blob/master/LICENSE) |
| `k6` | AGPL-3.0 | The load test. | v0.9.0 | [LICENSE](https://github.com/grafana/k6/blob/master/LICENSE.md) |

## Infrastructure Images

The local stack uses these images.
Each image carries its own licence.

| Image | Licence | Use in Kibtab | First version | Licence text |
| --- | --- | --- | --- | --- |
| `caddy:2-alpine` | Apache-2.0 | The edge proxy and the TLS. | v0.8.0 | [LICENSE](https://github.com/caddyserver/caddy/blob/master/LICENSE) |
| `postgres:16-alpine` | PostgreSQL License | The database in the local stack. | v0.8.0 | [COPYRIGHT](https://github.com/postgres/postgres/blob/master/COPYRIGHT) |

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