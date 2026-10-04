# Testing Kibtab

Version 0.0.0-docs. This document holds the rules for tests and for coverage.

Read [the agent rules](../AGENTS.md) section 4 for the rules that bind an agent.
Read [the contributing guide](../CONTRIBUTING.md) for the steps to send a change.

## The Coverage Aim

The module must reach 100% coverage.

* Cover every line and every branch of the code that you add or change.
* Cover the error path. A happy-path test alone is not a test.
* Cover the boundary. Test the empty case and the overflow case.
* Write a test for each bug before you write the fix.
* Never lower the floor to make a build pass.
* Never skip a test to make a build pass.
* Never use `//go:build ignore` or a blank identifier to hide a line.

The floor lives in `Makefile` as `COVERAGE_FLOOR`.

## Run The Coverage

```bash
make cover
make cover-verify
```

`make cover` prints the per-function breakdown.
`make cover-verify` fails below the floor.

Read the output before you send a change.
A file below the floor names the functions that need a test.

## The Layers Of Tests

Kibtab holds three kinds of test.

| Kind | Runs against | Needs a container |
| --- | --- | --- |
| Unit | The core. It uses a fake for each port. | No |
| Contract | One adapter. It holds the adapter to a port. | For an engine |
| Integration | Kibtab against a real database. | Yes |

* Keep the core test suite free of a container.
* Run the integration suite against a real database container.
* Run the contract suite for an adapter. Never repeat it in each adapter.
* Keep the client suite against a fake transport.

## The Rules For A Test

* Test one concern per test. Split a test that asserts two things.
* Name each test for the concern that it covers.
* Use a table-driven test for a set of related cases.
* Keep a fixture in one file. Share it.
* Keep the core test free of an adapter. Use a fake for a port.
* Cover the code that a port supplies. A fake covers the port, not the
  adapter.

## Contract Suites

A contract suite holds an adapter to a port.
It runs the same cases against every adapter.

* Write the suite once, against the port.
* Run it for each engine that ships.
* Run it for each client that ships.
* Add the suite to CI when the first adapter lands.

The `RowRepository` suite arrives at v0.3.0.
The `SyncService` suite arrives at v0.5.0.

## Coverage By Release

| Version | Aim |
| --- | --- |
| v0.2.0 | 100% on `internal/core/`. |
| v0.3.0 | 100% on the PostgreSQL adapter. |
| v1.0.0 | 100% on the whole module. |

## What Is Not Covered

The build does not measure these files.
Keep them thin so the gap stays small.

* The generated code. Never generate Go in this repository.
* The wiring in `cmd/kibtab/`. Cover it with one start-up test.

## Next Steps

* Read [the conventions guide](conventions.md) for the code rules.
* Read [the architecture guide](architecture.md) for the tree and the ports.
* Read [the contributing guide](../CONTRIBUTING.md) to send a change.