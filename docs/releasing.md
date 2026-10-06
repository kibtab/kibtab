# Releasing Kibtab

Version 0.0.0-docs. This document holds the versioning and the release steps.

Read section 5 of [the agent rules](https://github.com/kibtab/kibtab/blob/main/AGENTS.md) for the rules that bind an agent.
Read [the changelog format](changelogs/README.md) for the file layout.
Read [the plan](https://github.com/kibtab/kibtab/blob/main/plan.md) for the scope of each version.

## The Version Number

Kibtab uses semantic versioning. The number tells a user what changed.

| Part | Meaning |
| --- | --- |
| MAJOR | The public interface changed. A client must change. |
| MINOR | Kibtab added a feature. An old client keeps working. |
| PATCH | Kibtab fixed a defect. No interface changed. |

The public interface is the REST paths, the JSON field names, and the CLI
flags. A change to any of them needs a MAJOR version.

## The Gate For A Release

Every release passes each check below.

* `CGO_ENABLED=0 go build ./...` succeeds with no environment variable.
* `CGO_ENABLED=0 go vet ./...` reports no problem.
* `CGO_ENABLED=0 go test ./...` passes.
* `make cover-verify` passes.
* `make docs-check` and `make notice-check` report no problem.
* The module has no dependency that needs a C compiler.
* The code follows [the conventions guide](conventions.md).
* The commits follow [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md).

## The Gate For A Ship

* Every check in the release gate passes.
* `make release-check` passes.
* The archive builds for each target with `CGO_ENABLED=0`.
* `ldd` reports no dynamic library for the Linux binary.
* `docs/changelogs/vX.Y.Z.md` exists. It states the upgrade steps.
* `docs/index.md` lists the release.
* The tag exists. The tag never moves.

## The Release Files

* Release with GoReleaser. Read `.goreleaser.yaml` before you change it.
* Write one file at `docs/changelogs/vX.Y.Z.md` for each release.
* Copy [the changelog template](changelogs/TEMPLATE.md) for the file.
* Follow [the changelog format](changelogs/README.md) for the sections.
* Use the commit scope of the component. Read
  [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md) for the scope table.
* Update `docs/index.md` and `docs/CREDITS.md` with the release notes.
* Update `docs/THIRD_PARTY_NOTICES.md` when a dependency changes.

## Sign Off Each Commit

Sign off each commit with the Developer Certificate of Origin.

```bash
git commit -s -m "feat(db-postgres): add the batch delta validation"
```

Read [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md) for the full DCO section.

## Cut The Release

Run the gate first.

```bash
make check
make release-check
```

Build the archives with a dry run.

```bash
make snapshot
```

Tag the release.

```bash
make tag
```

The target reads the highest version from `docs/changelogs/v*.md`.
It warns when the changelog file for the tag is missing.
It never moves a published tag.

## After The Release

* Publish the announcement.
* Sign the release artefact.
* Update [the plan](https://github.com/kibtab/kibtab/blob/main/plan.md) when the next version opens.

## Next Steps

* Read [the changelog format](changelogs/README.md) for the file layout.
* Read [the upgrade guide](self-hosting/upgrade.md) for an upgrade.
* Read [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md) to send a change.