# Vendored Skill: simple-english

This directory holds a copy of the `simple-english` agent skill.

## Source

| Item | Value |
| --- | --- |
| Upstream repository | `https://github.com/AminBlg/SimpleEnglish` |
| Pinned commit | `32ea2d3f4404bfbff162c1861e5ae4f1bcc628b3` |
| Commit date | 2026-09-30 |
| Upstream skill version | `2.1.1` |
| Licence | MIT (see `LICENSE` in this directory) |
| Author | AminBlg |

The upstream project publishes the skill with no runtime dependency.

## What This Copy Holds

- `SKILL.md` holds the writing rules for agents.
- `references/rule-catalog.md` holds the 53 numbered rules of Issue 9.
- `references/word-swaps.md` holds the approved word substitutions.
- `references/strict-vocabulary.md` holds the dictionary rules for Strict mode.
- `references/use-cases.md` holds the task examples.

This copy holds the skill only. Kibtab does not keep the upstream tools.

## Why Kibtab Vendors It

The rules in `AGENTS.md` section 3 refer to this directory.
A vendored copy makes the rules stable.
Upstream changes must not alter the output of a Kibtab release without a
review.

This copy is the only source of the Simplified Technical English rules.
Kibtab does not keep the ASD-STE100 word lists.
The standard states that no reproduction of it is allowed without written
authority from ASD.

## How To Refresh The Copy

Run this command at a new upstream commit.
Then update the pinned commit above.
Record the change in `docs/changelogs/`.

```bash
git clone --depth 1 https://github.com/AminBlg/SimpleEnglish.git /tmp/simple-english
rm -rf skills/simple-english/references skills/simple-english/SKILL.md
cp -r /tmp/simple-english/skills/simple-english/. skills/simple-english/
cp /tmp/simple-english/LICENSE skills/simple-english/LICENSE
```

Read the upstream `CHANGELOG.md` before you copy.
Then confirm the rules in `AGENTS.md` still match `SKILL.md`.
Update `AGENTS.md` if they do not.