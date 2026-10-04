# Changelog Format

Version 0.0.0-docs. This folder holds one file for each Kibtab release.

The file name holds the version number.
The file for the release `v0.5.0` is `v0.5.0.md`.
Do not add a prefix, a suffix, or a date.
Do not create a file called `latest.md`.

No release file exists yet.

## The Format

Copy [the changelog template](TEMPLATE.md) to `vX.Y.Z.md` before you write it.
Use these headings in this order.

| Order | Heading | Content |
| --- | --- | --- |
| 1 | Title | `# Kibtab <version>` |
| 2 | Release line | The date and the commit |
| 3 | `Added` | The new features |
| 4 | `Changed` | The changes to a feature that works |
| 5 | `Fixed` | The defect repairs |
| 6 | `Removed` | The features that no longer work |
| 7 | `Security` | The security repairs. Write `None.` if there are none. |
| 8 | `Deprecated` | The features that a later version removes |
| 9 | `Upgrade` | The steps for a user |
| 10 | `Notes` | Anything else a user needs |

Remove a section when it has no entry.
The `Added` and `Upgrade` sections must stay.

## The Rules

* Create the file at the start of the release work. Do not create it at the
  tag.
* Put the version number in the first line of the file.
* Write in Simplified Technical English and British English.
  The rules are in `../../AGENTS.md` section 3.
* Keep each sentence under 20 words.
* Add the links to the issues and to the pull requests.
* The `Upgrade` section holds numbered steps. Put the condition before each
  command.
* Name each environment variable in the text.
* Update `../index.md` in the same commit as the file.
* Mark the matching boxes in `../../plan.md` when the code and the file both
  exist.

## The Version Rule

A change to the REST path, the JSON field name, or the CLI flag needs a MAJOR
version.
A new feature needs a MINOR version.
A defect repair needs a PATCH version.
The full rule is in `../../plan.md` section 1.

## The First File

The first file is `v0.1.0.md`.
Its `Added` section holds the scope of `../../plan.md` section 2.
Its `Released` line holds `Not released.` until the tag exists.