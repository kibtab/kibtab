# Runbook

Version 0.0.0-docs. This document holds the runbook for Kibtab operations.

Read [the deployment guide](self-hosting/deployment.md) for the first start.
Read [the upgrade guide](self-hosting/upgrade.md) to apply a change.
Read [the troubleshooting guide](self-hosting/troubleshooting.md) when a fault appears.

## Quick Start

### First Start

```bash
make tools
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
export KIBTAB_TABLES=my_table
docker compose up -d --build
```

### Verify Instance is Running

```bash
curl localhost:8080/healthz
# Expected: {"status":"ok","version":"devel"} or {"status":"ok","version":"v0.9.0"}
```

### Write a Row Through the Proxy

```bash
curl -X POST -H "Content-Type: application/json" \
  -H "X-Kibtab-User: admin" \
  -d '{"row":"row1","cells":[{"field":"k","value":"test"}]}' \
  https://localhost:8080/v1/tables/my_table/rows
```

### Check Metrics

```bash
curl http://localhost:8080/metrics
# Expected: kibtab_requests_total 1
```

### Check Audit

```bash
curl https://localhost:8080/v1/tables/my_table/rows/1/audit?field=k
```

## Daily Operations

### Retention Job

The audit retention job runs every 24 hours and removes audit rows older than
`KIBTAB_AUDIT_RETENTION_DAYS` (default 90).

```bash
# Check retention settings
curl https://localhost:8080/metrics | grep kibtab_audit

# Manually run retention
kibtab prune-audit
```

### Metrics Monitoring

Monitor these metrics Prometheus-style:

```
kibtab_requests_total - Total number of requests
kibtab_requests_by_route{route="/tables",status="200"} - Table listing requests
kibtab_requests_by_route{route="/v1/tables/{table}/rows",status="201"} - Row creation requests
kibtab_requests_by_status{status="200"} - Successful requests
kibtab_requests_by_status{status="429"} - Rate limited requests
kibtab_request_latency_seconds_sum - Total latency
kibtab_request_latency_seconds_count - Request count
kibtab_request_latency_seconds_bucket{le="0.1"} - Requests under 100ms
```

### Scale Up

To add more tables:

```bash
export KIBTAB_TABLES=my_table,other_table,another_table
docker compose up -d --build
# Verify
curl localhost:8080/healthz
```

## Troubleshooting

### Instance Won't Start

```bash
# Check if PostgreSQL is running
docker compose ps

# Check logs
docker compose logs kibtab

# Check DATABASE_URL
echo $DATABASE_URL
```

### Metrics Not Showing

```bash
# Verify /metrics endpoint
curl http://localhost:8080/metrics

# Check if metrics middleware is enabled
# (should be automatic in the binary)
```

### Audit Rows Not Appearing

```bash
# Check retention days
echo $KIBTAB_AUDIT_RETENTION_DAYS

# Manually prune old audit rows
kibtab prune-audit
```