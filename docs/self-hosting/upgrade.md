# Upgrading Kibtab

Version 0.0.0-docs. This document holds the steps to update a running install.

Read [the index](README.md) for the other guides.
Read [the backup guide](backup-restore.md) before you start.
Read [the troubleshooting guide](troubleshooting.md) when a step fails.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
No upgrade applies before that release.

## Before You Start

Read the changelog for your version.
Read [the changelog format](../changelogs/README.md) for the file layout.

```bash
ls docs/changelogs/
```

Take a backup before each upgrade.
Read [the backup guide](backup-restore.md) for the steps.

Check the new version before the change.

```bash
kibtab --version
```

## Upgrade With Docker

Pull the change and start the stack again.

```bash
git pull
docker compose up -d
```

The instance runs the migrations at the start.
Wait for each service to report the state.

```bash
docker compose ps
curl localhost:8080/healthz
```

## Upgrade A Binary

Replace the old binary with the new archive.
Check the checksum before you replace the file.

```bash
sha256sum --check kibtab_0.1.0_checksums.txt
```

Restart the service after the swap.

```bash
systemctl restart kibtab
curl localhost:8080/healthz
```

## Upgrade Across A Major Version

Read the migration guide for each step.
A change to a REST path, a JSON field name, or a CLI flag needs a MAJOR
version.

1. Read the migration guide for your version.
2. Take a backup.
3. Run the instance once with the new binary.
4. Rebuild the client with the taskpane of the new version.
5. Check one table in the database.

```bash
kibtab migrate --check
kibtab migrate --apply
```

## Roll Back

Stop the instance when a step fails.
Restore the data only when the migrations changed it.

```bash
docker compose down
git checkout <previous-tag>
docker compose up -d
```

The rollback does not undo a migration.
Read [the backup guide](backup-restore.md) for the restore step.

## Next Steps

* Read [the backup guide](backup-restore.md) to save the data.
* Read [the troubleshooting guide](troubleshooting.md) when a step fails.
* Read [the configuration guide](configuration.md) for each setting.