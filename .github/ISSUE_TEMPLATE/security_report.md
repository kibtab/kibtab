---
name: "Security Report"
about: "Report a security issue in private."
title: ""
labels: "security"
assignees: ""
---

<!--
Search the open and the closed issues before you create one.

===== READ CAREFULLY =====

This form is not the private route. It asks for the details before
you send them. Use it only when a private advisory does not fit.

For a private report, open an advisory at:
    https://github.com/kibtab/kibtab/security

Do not paste personal data into this form. Read SECURITY.md.
-->

# Security Report

## Describe The Security Issue

<!--
Give a clear and short description of the issue.
Name the trust boundary that the fault crosses.
-->

## Which Part Is Affected

<!--
Name the part that the fault touches.

- [ ] The client. A spreadsheet sends a value or a table name.
- [ ] The HTTP layer. A request carries an untrusted input.
- [ ] The core. A rule lets a bad value reach a port.
- [ ] A database engine. A query builds a statement from a name.
- [ ] The deployment. A setting opens a path to the data.
-->

## Steps To Reproduce

<!--
Give an unambiguous set of steps.
Keep real customer data out of the report.
Use a table that holds no personal data.
-->

1.
2.
3.

## Impact

<!--
Say what an attacker gains.
Say which table and which column the fault touches.
-->

## Your Environment

- Version used:
- Operating system:
- PostgreSQL version:
- Spreadsheet client and version: