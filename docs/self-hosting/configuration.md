# Configuring Kibtab

Version 0.0.0-docs. This document holds every setting for the Kibtab engine.

Read [the index](README.md) for the other guides.
Read [the deployment guide](deployment.md) for the first start.

## Status

Kibtab is pre-release.
The settings below arrive with each release.
Read [the plan](../../plan.md) for the version that adds each setting.

## The Variables

The engine reads each setting from the environment.

| Variable | Default | Purpose |
| --- | --- | --- |
| `PORT` | `8080` | The port for the API. |
| `DATABASE_URL` | None | The connection string for the database. |
| `KIBTAB_TABLES` | None | The tables in the registry. |
| `KIBTAB_DB_MAX_CONNS` | `20` | The size of the connection pool. |
| `KIBTAB_DB_QUERY_TIMEOUT` | `5s` | The timeout for one query. |
| `KIBTAB_AUDIT_RETENTION_DAYS` | `90` | The days to keep an audit row. |
| `KIBTAB_ENGINE` | `postgres` | The engine package that serves the data. |
| `KIBTAB_TRANSPORT` | `http` | The transport that serves the clients. |
| `KIBTAB_LOG_LEVEL` | `info` | The level of the log. |

The engine exits at start when a value is missing.
The engine exits at start when a chosen adapter has no package.

## The Database

Set the connection string before the first start.

```bash
DATABASE_URL=postgres://kibtab_user:PASSWORD@db:5432/kibtab_db?sslmode=disable
```

Keep the password in a secret store.
Do not commit the value to the repository.

## The Table Registry

Set each table in the registry. Use a comma between the names.

```bash
KIBTAB_TABLES=orders,customers,inventory
```

The registry limits the tables that a request can name.
The engine rejects a table that the registry does not hold.

## The Pool

Set the size of the connection pool for the load.

```bash
KIBTAB_DB_MAX_CONNS=20
KIBTAB_DB_QUERY_TIMEOUT=5s
```

The timeout must stay below the timeout of the proxy.

## The Audit

Set the days to keep an audit row.

```bash
KIBTAB_AUDIT_RETENTION_DAYS=90
```

The retention job removes old rows at the start of each day.

## Check The Settings

The route `GET /healthz` returns the version.
The route `GET /metrics` returns the pool counters.

```bash
curl localhost:8080/metrics
```

## Next Steps

* Read [the deployment guide](deployment.md) for the first start.
* Read [the upgrade guide](upgrade.md) to apply a change.
* Read [the troubleshooting guide](troubleshooting.md) when a fault appears.