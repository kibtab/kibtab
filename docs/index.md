# Kibtab Documentation

Version 0.0.0-docs. This page lists the Kibtab documents.

Kibtab turns a spreadsheet into a client for a relational database.
The engine holds a version for each table.
It rejects a write from an old client.

## Read The Documents

| Document | Purpose |
| --- | --- |
| [Readme](../README.md) | The project summary. It is the front page. |
| [Plan](../plan.md) | The release plan. It splits the work by version. |
| [Agent Rules](../AGENTS.md) | The rules for AI coding agents. |
| [Install](install.md) | The install steps for a binary or a source build. |
| [Self-hosting](self-hosting/README.md) | The section for running Kibtab on a server. |
| [Architecture](architecture.md) | The layers, the ports, and the extension steps. |
| [Licence](licence.md) | The terms for the engine and the skill. |
| [Third Party Notices](THIRD_PARTY_NOTICES.md) | The licences of the dependencies. |
| [Credits](CREDITS.md) | The people, the projects, and the sources. |
| [Changelog rules](changelogs/README.md) | The format for each release file. |
| [Contributing](../CONTRIBUTING.md) | The steps to send a change. |
| [Code of Conduct](../CODE_OF_CONDUCT.md) | The terms for taking part. |
| [Security](../SECURITY.md) | The private route for a vulnerability. |

The plan in `plan.md` holds the scope for each release.
Read it before you start a task.

## Releases

Each release has one file in the [changelogs](changelogs/) folder.
The file name holds the version number.
No release file exists yet.
The plan in `plan.md` holds the scope for each version.

| Version | Theme | Changelog file |
| --- | --- | --- |
| v0.1.0 | Skeleton | Not written yet |
| v0.2.0 | Domain and ports | Not written yet |
| v0.3.0 | PostgreSQL adapter | Not written yet |
| v0.4.0 | Write path | Not written yet |
| v0.5.0 | HTTP API | Not written yet |
| v0.6.0 | Office.js client | Not written yet |
| v0.7.0 | Audit and locking | Not written yet |
| v0.8.0 | Packaging | Not written yet |
| v0.9.0 | Hardening | Not written yet |
| v1.0.0 | Stable | Not written yet |
| v1.1.0 to v1.9.0 | Growth | Not planned yet |
| v2.0.0 and later | Breaking interface change | Not planned yet |

Copy `changelogs/TEMPLATE.md` to `changelogs/vX.Y.Z.md` at the start of the
release work.
Do not write a file for a planned release that has no code.

## Documents That A Release Adds

The documents below exist now as pre-release drafts.
Each one updates when its release arrives.

| Document | Added in | Purpose |
| --- | --- | --- |
| `README.md` | v0.1.0 | The project summary. |
| `CONTRIBUTING.md` | v0.1.0 | The steps to send a change. |
| `CODE_OF_CONDUCT.md` | v0.1.0 | The terms for taking part. |
| `SECURITY.md` | v0.1.0 | The private route for a vulnerability. |
| `.github/ISSUE_TEMPLATE/` | v0.1.0 | The forms for a bug, a feature, and a security report. |
| `docs/install.md` | v0.1.0 | The install steps. |
| `docs/architecture.md` | v0.2.0 | The layers and the port contract. |
| `docs/licence.md` | v0.1.0 | The terms for the engine and the skill. |
| `docs/openapi.yaml` | v0.5.0 | The REST interface definition. |
| `docs/self-hosting/` | v0.8.0 | The section for running Kibtab on a server. |
| `docs/runbook.md` | v0.9.0 | The steps for an operator. |
| `docs/licences/` | v0.8.0 | A copy of each dependency licence. |

## Documentation Rules

Write each document in two styles at the same time.
The first style is Simplified Technical English.
The second style is British English.
The full rules are in `../AGENTS.md` section 3.

* Load the vendored skill first. Read `../skills/simple-english/SKILL.md`.
* Keep a sentence under 20 words.
* Use active voice and simple tenses.
* Define each technical term at its first use.
* Do not use em dashes, contractions, `should`, `would`, `may`, or `might`.
* Use the approved words in `../skills/simple-english/references/`.
* Use British English spellings. The forms are `licence`, `artefact`, and
  `behaviour`.

Run the check before you open a pull request.

```bash
make docs-check
```

The target checks the sentence length and the banned words.
The target does not check the word lists.
The vendored skill is the source of the word rules.
Kibtab does not keep a copy of the ASD-STE100 word lists.
The standard states that no reproduction of it is allowed without written
authority from ASD.

## The Vendored Skill

The skill lives in `../skills/simple-english/`.
Read `../skills/simple-english/VENDOR.md` for the pinned commit.
Read `../skills/simple-english/LICENSE` for the terms.

## Which Guide Owns Which Step

Each guide owns one task. Do not copy a step into a second file.

| Guide | Owns |
| --- | --- |
| [install.md](install.md) | The install and build steps. |
| [self-hosting/README.md](self-hosting/README.md) | The index for the self-hosting section. |
| [self-hosting/deployment.md](self-hosting/deployment.md) | The steps to deploy the stack on a server. |
| [self-hosting/configuration.md](self-hosting/configuration.md) | Each setting and each environment variable. |
| [self-hosting/upgrade.md](self-hosting/upgrade.md) | The steps to update a running install. |
| [self-hosting/backup-restore.md](self-hosting/backup-restore.md) | The steps to save and to restore the data. |
| [self-hosting/troubleshooting.md](self-hosting/troubleshooting.md) | The steps to diagnose a fault. |
| [architecture.md](architecture.md) | The layers, the ports, and the extension steps. |
| [licence.md](licence.md) | The licence terms and the obligations. |

The rules for this file are in the section above.

## Build And Toolchain Documents

| File | Purpose |
| --- | --- |
| `air.toml` | The hot reload settings for a local run. |
| `.goreleaser.yaml` | The release build settings. |
| `docker-compose.yml` | The local stack with Caddy, the engine, and PostgreSQL. |
| `Dockerfile` | The multi-stage build. It uses a `scratch` base. |
| `Makefile` | The targets `build`, `test`, `lint`, `run`, `docs-check`, `notice-check`, and `release-check`. |

Kibtab sets `CGO_ENABLED=0` in every build.
No dependency needs a C compiler.