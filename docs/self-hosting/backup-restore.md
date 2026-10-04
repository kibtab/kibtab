# Backup And Restore

Version 0.0.0-docs. This document holds the steps to save and to restore the
data of Kibtab.

Read [the index](README.md) for the other guides.
Read [the upgrade guide](upgrade.md) before you restore.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
The steps below apply from v0.1.0.

## What To Save

The database holds all of the data.
The `_kibtab_meta` schema holds the versions and the audit rows.
Save the whole database. Do not save one table.

## Take A Backup

Run the command on a schedule.
Write the file outside the container.

```bash
docker compose exec db pg_dump -U kibtab_user kibtab_db > backup.sql
```

Compress the file after the dump.

```bash
gzip backup.sql
```

## Check The Backup

Read the file before you trust it.
A backup that was never read is not a backup.

```bash
gzip -t backup.sql.gz
```

## Restore

Stop the instance before a restore.
Keep the API closed while the data changes.

```bash
docker compose stop engine
```

Load the file into the database.

```bash
docker compose exec -T db psql -U kibtab_user kibtab_db < backup.sql
```

Start the instance after the load.

```bash
docker compose start engine
curl localhost:8080/healthz
```

## The Audit Rows

Set the retention days before the retention job runs.
Read [the configuration guide](configuration.md) for the variable.

```bash
KIBTAB_AUDIT_RETENTION_DAYS=90
```

The job removes the rows at the start of each day.
A restore brings back the rows that the job removed.
Set the retention days again after a restore.

## Keep The Secrets

The dump holds the data of every table.
It holds no password.
Keep the file in a secret store.
Do not commit the file to the repository.

## Next Steps

* Read [the upgrade guide](upgrade.md) to apply a new version.
* Read [the troubleshooting guide](troubleshooting.md) when a load fails.
* Read [the configuration guide](configuration.md) for each setting.