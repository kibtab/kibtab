# STE100 Technical Names: tooling terms

Version 0.0.0-docs. This document holds the tooling Technical Names.

These terms name the commands, targets, files, and gates of Kibtab.
Read the [STE100 Technical Names](index.md) dictionary first.

## Technical Name: Target
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** A named task in the `Makefile` that runs one command.
- **Approved Form:** Target (singular), Targets (plural)
- **Do Not Use:** Rule, Task, Goal, Step
- **Correct Example:** *The **target** `make check` runs every gate.*
- **Incorrect Example:** *The rule `make check` runs every gate.*

## Technical Name: Gate
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** A check that stops a build when it fails.
- **Approved Form:** Gate (singular), Gates (plural)
- **Do Not Use:** Check, Fence, Barrier, Filter
- **Correct Example:** *The release **gate** must pass before a tag.*
- **Incorrect Example:** *The release check must pass before a tag.*

## Technical Name: Floor
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** The lowest coverage percentage that a build accepts.
- **Approved Form:** Floor (singular), Floors (plural)
- **Do Not Use:** Threshold, Limit, Minimum, Baseline
- **Correct Example:** *The coverage **floor** is 100 per cent.*
- **Incorrect Example:** *The coverage threshold is 100 per cent.*

## Technical Name: Artefact
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** A file that a build produces. It holds no source code.
- **Approved Form:** Artefact (singular), Artefacts (plural)
- **Do Not Use:** Artifact, Output, Product, Result
- **Correct Example:** *The release **artefact** holds the binary and the licence.*
- **Incorrect Example:** *The release artifact holds the binary and the license.*

## Technical Name: Toolchain
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** The set of tools that build and test the Go sources.
- **Approved Form:** Toolchain (singular), Toolchains (plural)
- **Do Not Use:** Tool Set, Suite, Environment, Stack
- **Correct Example:** *The **toolchain** needs Go 1.22 and no C compiler.*
- **Incorrect Example:** *The tool set needs Go 1.22 and no C compiler.*

## Technical Name: Vendored
- **Part of Speech:** Modifier
- **Category:** Tooling Term
- **Definition:** Describes a copy of another project that Kibtab stores.
- **Approved Form:** Vendored
- **Do Not Use:** Embedded, Bundled, Frozen, Pinned Copy
- **Correct Example:** *The **vendored** skill stays at its pinned commit.*
- **Incorrect Example:** *The embedded skill stays at its pinned commit.*

## Technical Name: Milestone
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** One version of Kibtab with one theme and one result.
- **Approved Form:** Milestone (singular), Milestones (plural)
- **Do Not Use:** Release, Stage, Phase, Sprint
- **Correct Example:** *The v0.1.0 **milestone** builds the skeleton.*
- **Incorrect Example:** *The v0.1.0 phase builds the skeleton.*

## Technical Name: Snapshot
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** A test build that GoReleaser makes without an upload.
- **Approved Form:** Snapshot (singular), Snapshots (plural)
- **Do Not Use:** Dry Run, Trial, Preview Build, Candidate
- **Correct Example:** *The `snapshot` **target** builds each **snapshot** archive.*
- **Incorrect Example:** *The `snapshot` target builds each preview build archive.*

## Technical Name: Canonical
- **Part of Speech:** Modifier
- **Category:** Tooling Term
- **Definition:** Describes the one file that holds a value. Change it first.
- **Approved Form:** Canonical
- **Do Not Use:** Authoritative, Official, Master, Primary
- **Correct Example:** *The **canonical** tree lives in `docs/architecture.md`.*
- **Incorrect Example:** *The authoritative tree lives in `docs/architecture.md`.*

## Technical Name: Baseline Version
- **Part of Speech:** Noun
- **Category:** Tooling Term
- **Definition:** The tool version that Kibtab requires. Kibtab pins it.
- **Approved Form:** Baseline Version (singular), Baseline Versions (plural)
- **Do Not Use:** Minimum Version, Required Version, Floor Version
- **Correct Example:** *The **baseline version** of Go is 1.22.*
- **Incorrect Example:** *The minimum version of Go is 1.22.*

## Next Steps

* Read [identifiers.md](identifiers.md) for the source names.
* Read [architecture-terms.md](architecture-terms.md) for the layer terms.