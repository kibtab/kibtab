# Installing Kibtab

Version 0.0.0-docs. This document holds the install steps for Kibtab.

Read [the self-hosting section](self-hosting/README.md) for the server stack.
Read [the architecture guide](architecture.md) for the design.

## Status

Kibtab is pre-release.
The first release is v0.1.0.
No binary is published yet.
The commands below apply from v0.1.0.

## Install With Go

Install the engine with the Go tool.

```bash
go install github.com/kibtab/kibtab/cmd/kibtab@latest
```

The command puts the binary in `$(go env GOPATH)/bin`.
Add that path to `PATH` before you run `kibtab`.

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
kibtab --version
```

## Install A Release Archive

Each release holds an archive for each target.
Use the archive when you cannot install Go.

1. Read the release page for your version.
2. Download the archive for your operating system.
3. Check the checksum against the release file.
4. Put the binary on your `PATH`.

```bash
sha256sum --check kibtab_0.1.0_checksums.txt
```

The archives use `tar.gz` for Linux and macOS.
The archives use `zip` for Windows.

## Build From Source

Clone the repository and build the engine.

```bash
git clone https://github.com/kibtab/kibtab.git
cd kibtab
CGO_ENABLED=0 go build ./cmd/kibtab
```

Kibtab needs no C compiler.
Every build sets `CGO_ENABLED=0`.
Do not unset that variable.

The build needs Go 1.22 or a newer version.

```bash
go version
```

## Run In Watch Mode

Use `air` during development. It rebuilds the engine on a code change.

```bash
go install github.com/air-verse/air@latest
air
```

Read `air.toml` before you change it.

## Run The Tests

```bash
CGO_ENABLED=0 go test ./...
```

The core tests need no database.
The engine tests need a PostgreSQL container.

## Check The Install

```bash
kibtab --version
curl localhost:8080/healthz
```

The health route returns the status and the version.

## Next Steps

* Read [the self-hosting section](self-hosting/README.md) for the server stack.
* Read [the architecture guide](architecture.md) for the code layout.
* Read [the plan](../plan.md) for the scope of each version.
* Read [the licence](licence.md) for the terms.