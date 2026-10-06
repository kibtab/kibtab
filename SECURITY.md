# Reporting Security Issues

Kibtab holds the data of a production database.
Kibtab writes to a table that a user trusts.
The maintainers take a security issue seriously.

Report a vulnerability in private.
Do not open a public issue.

Read [the documentation index](https://kibtab.readthedocs.io/en/latest/) for every guide.
Read [the agent rules](AGENTS.md) before you change code for a report.

## How To Report

Use the private route that matches the urgency.

* Open a private security advisory at
  [github.com/kibtab/kibtab/security](https://github.com/kibtab/kibtab/security).
* Open a GitHub issue with the label `security` when no advisory fits.

A report helps the maintainers most when it gives these items.

* The version of Kibtab.
* The steps to reproduce the fault.
* The table and the field that the fault touches.
* The impact on the data.

Keep a proof of concept small.
Keep real customer data out of a report.
Use a table that holds no personal data.

## What We Look At First

Kibtab writes to a table from a spreadsheet.
The threat model for that path holds these items.

* A table name or a field name that comes from a request.
* A value that the client sends without a check.
* A cell write that skips the version check.
* A query that a client can change through a parameter.
* A spreadsheet that a user opens from an unknown sender.

The **instance** rejects a table that the registry does not hold.
The **instance** compares the version before it writes.

## What A Fix Must Do

* Keep the version check in the same transaction as the write.
* Use a parameterized query for every value.
* Quote every identifier through the `Dialect` port.
* Keep a secret out of a log line and out of an error message.
* Keep the default of Caddy strict for a rate limit.

## What Is Out Of Scope

These reports do not qualify.

* A defect with no path from an untrusted input.
* A missing header on a deployment that the operator controls.
* A brute force attempt against a deployment with no rate limit.
* A report from an automated scan with no proof.

Read [the architecture guide](https://kibtab.readthedocs.io/en/latest/architecture.html) for the trust boundary.
Read [the conventions guide](https://kibtab.readthedocs.io/en/latest/conventions.html) for the code rules.
Read [the licence guide](https://kibtab.readthedocs.io/en/latest/licence.html) for the terms.

## After A Report

The maintainers acknowledge each report.
They give an estimate for a fix.
They credit a report in the release notes when you ask for it.

Report a vulnerability in private.
Thank you for the report.