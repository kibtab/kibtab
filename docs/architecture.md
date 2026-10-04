# Architecture

Version 0.0.0-docs. This document holds the design of Kibtab.
It holds the layer rules and the port contract.

This document holds the canonical layout and the canonical port table.
Read [AGENTS.md](../AGENTS.md) section 1 for the rules for agents.

## The Goal

Kibtab is a spreadsheet client and a database server.
It must serve any relational database.
It must serve any spreadsheet.

The core must not name one engine.
The core must not name one client.
A new engine needs one package.
A new client needs one folder.

## The Repository Tree

This tree is canonical. Change it before you add a file or a folder.

```text
kibtab/
├── cmd/kibtab/              The wiring. It builds each adapter.
├── internal/
│   ├── core/                The kernel. It has no I/O.
│   │   ├── domain/          The models.
│   │   ├── ports/           The interfaces. The core owns them.
│   │   └── services/        The use cases.
│   └── adapters/
│       ├── driven/<engine>/ One package per database engine.
│       └── driver/<transport>/  One package per transport.
├── client/<spreadsheet>/    One folder per spreadsheet client.
├── docs/                    The documentation and the changelogs.
│   ├── changelogs/          One file per release.
│   └── self-hosting/        The guides for a server.
├── scripts/                 The checks that run without a container.
├── skills/simple-english/   The vendored writing skill.
├── AGENTS.md                The rules for agents.
├── CONTRIBUTING.md          The steps to send a change.
├── plan.md                  The implementation plan.
├── Makefile                 The build, test, and release targets.
├── air.toml                 The hot reload settings.
├── Caddyfile                The edge proxy and the TLS.
├── docker-compose.yml       The local stack.
├── Dockerfile               The multi-stage build.
├── .goreleaser.yaml         The release build.
├── .gitignore               The build and coverage output.
└── go.mod                   The module. It sets Go 1.22.
```

A file in `internal/core/` never imports a package under `internal/adapters/`.
A file under `internal/adapters/` never imports another adapter.
The core imports nothing outside the standard library.

## The Ports

The core names a need. The adapter supplies the answer.
The table lives in this document.

| Port | Need it fills |
| --- | --- |
| `RowRepository` | Read and write rows in a table. |
| `TableRegistry` | List the tables. Read the metadata for a table. |
| `AuditWriter` | Record a change to a cell. |
| `Dialect` | Quote an identifier. Map a value to a type. Page a query. |
| `TransactionRunner` | Run a group of writes as one unit. |
| `Clock` | Give the current time. |
| `SyncService` | Accept a sync payload. Return a sync result. |

Add a port to the table in this document first.
Then define it in `internal/core/ports/`.

## The Data Flow

A request moves through the layers in this order.

1. A transport package under `driver/` reads the request.
2. The transport maps the request to a domain value.
3. A service in `internal/core/services/` calls the ports.
4. An engine package under `driven/` runs the query.
5. The engine maps the result to a domain value.
6. The transport maps the domain value to a response.

The core runs at steps 3 and 4.
The core holds no SQL and no JSON.

## Add A Database

Add an engine without a change to the core.

1. Create `internal/adapters/driven/<engine>/`.
2. Implement the `Dialect` port for that engine.
3. Implement `RowRepository` for that engine.
4. Implement `TransactionRunner` for that engine.
5. Implement `TableRegistry` for that engine.
6. Add a constructor in `cmd/kibtab/`.
7. Read the engine name from the environment.
8. Run the contract suite for that engine.

The core needs no change.

## Add A Spreadsheet

Add a client without a change to the core.

1. Create `client/<spreadsheet>/`.
2. Send the sync payload to a transport under `driver/`.
3. Map each value of the spreadsheet to a domain value.
4. Add the manifest for the spreadsheet.
5. Run the contract suite for `SyncService`.

The core needs no change.

## Add A Transport

1. Create `internal/adapters/driver/<transport>/`.
2. Map the request to a domain value.
3. Call the service in `internal/core/services/`.
4. Map the domain value to a response.
5. Run the contract suite for `SyncService`.

The core needs no change.

## The Tests

A contract test holds each adapter to the port.

* The core test suite needs no database.
* The engine suite runs against a real database container.
* The client suite runs against a fake transport.
* The CI step fails the build when the core imports an adapter.
* The CI step fails the build when the core names an engine or a client.

## Next Steps

* Read [the plan](../plan.md) for the version that adds each adapter.
* Read [AGENTS.md](../AGENTS.md) section 1 for the rules.
* Read [the install guide](install.md) to build the engine.