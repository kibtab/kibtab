# Testing Kibtab

Version 0.0.0-docs. This document holds the rules for tests and for coverage.

Read section 4 of [the agent rules](https://github.com/kibtab/kibtab/blob/main/AGENTS.md) for the rules that bind an agent.
Read [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md) for the steps to send a change.

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

## PostgreSQL Testing

The integration tests need a real PostgreSQL. The Makefile manages this
for you with dedicated targets.

### Start A Testing Instance

```bash
make postgres-start
```

This target starts a dedicated PostgreSQL instance for testing:

* It uses port 5433 by default to avoid conflicts with existing PostgreSQL.
* It creates a test user `kibtab_test` and database `kibtab_test`.
* It sets `KIBTAB_TEST_DATABASE_URL` for the test run.
* It works with a local `pg_ctl` installation or Docker.

### Stop The Testing Instance

```bash
make postgres-stop
```

This target stops the testing instance and cleans up:

* It stops the PostgreSQL process or Docker container.
* It removes temporary data directories.
* It does not affect any existing PostgreSQL installation.

### Check Instance Status

```bash
make postgres-status
```

This target reports whether the testing instance is running.

### Run Coverage With PostgreSQL

```bash
make cover-verify
```

`make cover-verify` now depends on `postgres-start`. It:

1. Starts a dedicated PostgreSQL instance.
2. Runs the tests with coverage.
3. Checks the coverage against the floor.
4. Stops the PostgreSQL instance.

The gate no longer skips when there is no database. It starts one.

### Manual Control

For advanced use, you can control the instance manually:

```bash
make postgres-start
make test
make cover-verify
make postgres-stop
```

Set the port and credentials with environment variables:

```bash
make PG_TEST_PORT=5434 PG_TEST_USER=myuser postgres-start
```

### Fallback Behaviour

When PostgreSQL is not installed and Docker is not available, the
`postgres-start` target prints a warning and continues. The coverage
gate then fails because the integration tests skip.

Install PostgreSQL or Docker to run full coverage checks locally.

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
| v0.5.0 | 100% on the whole module. |

## What Is Not Covered

The build does not measure these files.
Keep them thin so the gap stays small.

* The generated code. Never generate Go in this repository.
* The wiring in `cmd/kibtab/`. Cover it with one start-up test.

## Next Steps

* Read [the conventions guide](conventions.md) for the code rules.
* Read [the architecture guide](architecture.md) for the tree and the ports.
* Read [the contributing guide](https://github.com/kibtab/kibtab/blob/main/CONTRIBUTING.md) to send a change.