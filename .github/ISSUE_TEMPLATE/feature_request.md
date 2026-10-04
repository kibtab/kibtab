---
name: "Feature Request"
about: "Suggest an idea for the project."
title: ""
labels: "feature"
assignees: ""
---

# Feature Request

Your idea could already exist.
Check the [open issues](https://github.com/kibtab/kibtab/issues) and the
[plan](https://github.com/kibtab/kibtab/blob/main/plan.md) before you create
one.

## Is Your Idea In The Plan

<!--
The plan holds the scope of each version.
Name the version that holds the idea, or write none.

    v0.4.0 - Write path

Read the plan before you write this report.
-->

## What Problem Does It Solve

<!--
Give a clear and short description of the problem.
Say what a user cannot do today.
-->

## Describe The Solution

<!--
Give a clear and short description of what you want.
-->

## Which Adapter Does It Touch

<!--
Name the part that the change touches.

- [ ] The core. The change reaches `internal/core/`.
- [ ] An engine. The change adds an engine package.
- [ ] A transport. The change adds a package under `driver/`.
- [ ] A client. The change adds a folder under `client/`.
- [ ] The deployment. The change touches the stack.
- [ ] None. The change touches the docs only.
-->

Keep the core free of an engine name and of a client name.
A change that adds an adapter must not change a file under `internal/core/`.

## Alternatives You Considered

<!--
Give a clear and short description of the other routes you considered.
-->

## Breaking Change

<!--
Does the change alter a REST path, a JSON field name, or a CLI flag?
A change to any of them needs a MAJOR version.
-->

## Additional Context

<!--
Add any other detail that helps.
-->

<!--
For a new spreadsheet or a new database, read docs/architecture.md first.
That guide holds the steps to add each adapter.
-->