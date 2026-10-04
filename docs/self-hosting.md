# Self-Hosting Kibtab

Version 0.0.0-docs. This document holds the steps to run Kibtab on your own
server.

Read [the install guide](install.md) for a binary install.
Read [the licence](licence.md) for the terms.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
The stack arrives at v0.8.0.
The file `docker-compose.yml` does not exist yet.

## The Stack

The stack holds three services.

| Service | Purpose |
| --- | --- |
| `caddy` | The edge proxy. It manages the TLS. |
| `engine` | The Kibtab engine. It serves the API. |
| `db` | PostgreSQL. It holds the data. |

Caddy handles the certificate.
The engine holds no TLS code.

## Requirements

* A server with a Linux host.
* A domain name with a DNS record for the server.
* Docker and the Docker Compose plugin.

Check the tools before you start.

```bash
docker --version
docker compose version
```

## Start The Stack

Clone the repository and start the stack.

```bash
git clone https://github.com/kibtab/kibtab.git
cd kibtab
docker compose up -d
```

Check each service after the start.

```bash
docker compose ps
curl localhost:8080/healthz
```

## Set The Domain

Caddy needs a domain name for the certificate.
Read the file `Caddyfile` before you start the stack.
Set the domain in that file.

```bash
docker compose down
docker compose up -d
```

Caddy requests a certificate from the ACME service at the first start.
The rate limit for that service is low.
Use a test domain first.

## Configuration

The engine reads each setting from the environment.

| Variable | Purpose |
| --- | --- |
| `PORT` | The port for the API. The default is `8080`. |
| `DATABASE_URL` | The connection string for PostgreSQL. |
| `KIBTAB_TABLES` | The tables in the registry. |
| `KIBTAB_DB_MAX_CONNS` | The size of the connection pool. |
| `KIBTAB_DB_QUERY_TIMEOUT` | The timeout for one query. |
| `KIBTAB_AUDIT_RETENTION_DAYS` | The days to keep an audit row. |

The variable `DATABASE_URL` holds the password.
Keep it in a secret store.
Do not commit it to the repository.

```bash
docker compose up -d
```

## The Health Check

The route `GET /healthz` returns the status and the version.

```bash
curl localhost:8080/healthz
```

Read the route `GET /metrics` for the counters.
The route needs no data from a user.

## Backup

Back up the database with `pg_dump`.
Run the command on a schedule.

```bash
docker compose exec db pg_dump -U kibtab_user kibtab_db > backup.sql
```

Test the restore before you need it.

## Update

Pull the change and start the stack again.

```bash
git pull
docker compose up -d
```

The engine runs the migrations at the start.
The runbook arrives at v0.9.0. Read `docs/runbook.md` when an update fails.

## Troubleshoot

Check the state of each service.

```bash
docker compose logs engine
docker compose logs caddy
```

The engine exits when it cannot reach the database.
The Caddy service exits when the domain name fails the check.

## Next Steps

* Read [the install guide](install.md) for a binary install.
* Read [the architecture guide](architecture.md) for the design.
* Read [the plan](../plan.md) for the version that adds the stack.