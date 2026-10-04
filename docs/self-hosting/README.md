# Self-Hosting Kibtab

Version 0.0.0-docs. This folder holds the guides for running Kibtab on a
server.

Read the guide for your task.
Do not follow a step from a second guide.
Each guide owns its own task.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
The stack arrives at v0.8.0.
The file `docker-compose.yml` does not exist yet.
Each guide in this folder is a pre-release draft.

## The Guides

| Guide | Owns |
| --- | --- |
| [deployment.md](deployment.md) | The steps to deploy the stack on a server. |
| [configuration.md](configuration.md) | Each setting and each environment variable. |
| [upgrade.md](upgrade.md) | The steps to update a running install. |
| [backup-restore.md](backup-restore.md) | The steps to save and to restore the data. |
| [troubleshooting.md](troubleshooting.md) | The steps to diagnose a fault. |

Read [the install guide](../install.md) for a binary install.
Read [the architecture guide](../architecture.md) for the design.
The runbook arrives at v0.9.0. Read `docs/runbook.md` after that release.

## The Stack

The stack holds three services.

| Service | Purpose |
| --- | --- |
| `caddy` | The edge proxy. It manages the TLS. |
| `engine` | The Kibtab engine. It serves the API. |
| `db` | PostgreSQL. It holds the data. |

Caddy handles the certificate.
The engine holds no TLS code.

## Choose A Path

Pick one path before you start.

* Run the stack with Docker. Read [deployment.md](deployment.md).
* Run a release archive on a bare host. Read [deployment.md](deployment.md).
* Update a running install. Read [upgrade.md](upgrade.md).

## Rules For This Section

* Write each guide for one task.
* Do not copy a step into a second guide.
* Add this index in the same commit as a new guide.
* Link each guide from this index. Leave no guide unlinked.