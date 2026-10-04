# Troubleshooting Kibtab

Version 0.0.0-docs. This document holds the steps to diagnose a fault.

Read [the index](README.md) for the other guides.
Read [the configuration guide](configuration.md) for each setting.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
The steps below apply from v0.1.0.

## Check The State

Read the state of each service first.

```bash
docker compose ps
```

Read the log of the service that failed.

```bash
docker compose logs engine
docker compose logs db
docker compose logs caddy
```

## The Engine Does Not Start

The engine exits when it cannot reach the database.
Check the value of `DATABASE_URL`.
Read [the configuration guide](configuration.md) for the variable.

```bash
docker compose logs engine
```

The engine exits when a chosen adapter has no package.
Check the value of `KIBTAB_ENGINE`.

## The Route Returns A Conflict

The engine rejects a write when the version does not match.
The sheet holds an old version.
Reload the row in the client.

Read [the architecture guide](../architecture.md) for the version check.

## The Query Times Out

The timeout of the engine sits above the timeout of the proxy.
Set the timeout of the engine below the timeout of the proxy.

```bash
KIBTAB_DB_QUERY_TIMEOUT=5s
```

The pool runs out when the load holds each connection.
Raise the size of the pool.

```bash
KIBTAB_DB_MAX_CONNS=40
```

## Caddy Holds No Certificate

Caddy needs a domain name with a DNS record for the server.
Read the file `Caddyfile` and check the domain.
The ACME service holds a low rate limit.
Use a test domain first.

```bash
docker compose logs caddy
```

## The Routes Fail After An Upgrade

A change to a REST path, a JSON field name, or a CLI flag needs a MAJOR
version.
Read the migration guide for your version.

```bash
kibtab --version
```

## Check The Metrics

The route `GET /metrics` returns the counters.

```bash
curl localhost:8080/metrics
```

The counters show the query count and the error count.

## Next Steps

* Read [the configuration guide](configuration.md) for each setting.
* Read [the deployment guide](deployment.md) for the first start.
* Read [the upgrade guide](upgrade.md) after a version change.