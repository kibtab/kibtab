# Deploying Kibtab

Version 0.0.0-docs. This document holds the steps to deploy Kibtab on a
server.

Read [the index](README.md) for the other guides.
Read [the configuration guide](configuration.md) for each setting.
Read [the troubleshooting guide](troubleshooting.md) when a service fails.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
The stack arrives at v0.8.0.
The file `docker-compose.yml` does not exist yet.
The steps below apply from v0.8.0.

## The Requirements

* A server with a Linux host.
* A domain name with a DNS record for the server.
* Docker and the Docker Compose plugin.

Check each tool before you start.

```bash
docker --version
docker compose version
```

## Deploy With Docker

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

The route `GET /healthz` returns the status and the version.

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

## Deploy A Binary

Deploy a release archive when the host has no Docker.

1. Read the release page for your version.
2. Install the archive for your operating system.
3. Put the binary on the `PATH`.
4. Run the instance under a service manager.

Read [the install guide](../install.md) for the install steps.

The instance needs the variables in
[the configuration guide](configuration.md).
Keep the password in a secret store.
Do not put the password in a unit file on disk.

## Deploy Behind A Proxy

Put a reverse proxy in front of the instance when the host has no Caddy.
The proxy must pass the WebSocket upgrade for the live sync route.
The proxy must set a timeout above the query timeout.

Read [the configuration guide](configuration.md) for the query timeout.

## Check The Deploy

```bash
curl localhost:8080/healthz
curl localhost:8080/metrics
```

The metrics route returns the counters and no data from a user.

## Next Steps

* Read [the configuration guide](configuration.md) for each setting.
* Read [the upgrade guide](upgrade.md) to update the install.
* Read [the backup guide](backup-restore.md) to save the data.
* Read [the troubleshooting guide](troubleshooting.md) when a fault appears.