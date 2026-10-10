#!/usr/bin/env python3
"""Compare go.mod with the dependency table in docs/THIRD_PARTY_NOTICES.md.

The script parses the require blocks of go.mod and the Go Dependencies
table of the notices. It fails when a module misses a row, when a row
names no release, or when a row has too few cells.

The notices table names the release that first ships each dependency. The
`go.mod` file holds the exact module version, so the script does not
compare module versions.

Usage:

    python3 scripts/notice-check.py
"""

import re
import sys

GO_MOD = "go.mod"
NOTICES = "docs/THIRD_PARTY_NOTICES.md"

# A Kibtab release version, such as v0.3.0.
RELEASE = re.compile(r"^v\d+\.\d+\.\d+$")

# One require line of go.mod, without the indirect marker.
REQUIRE = re.compile(r"^\s+([A-Za-z0-9._/-]+)\s+(v\d[^ \t]*)\s*//\s*indirect")
REQUIRE_DIRECT = re.compile(r"^\s+([A-Za-z0-9._/-]+)\s+(v\d[^ \t]*)\s*$")

# One table row of the notices. The first cell names the module in backticks.
ROW = re.compile(r"^\|\s*`([^`]+)`\s*\|")


def read_go_mod():
    """Return the module versions of go.mod as {module: version}."""
    modules = {}
    with open(GO_MOD, encoding="utf-8") as handle:
        for line in handle:
            match = REQUIRE_DIRECT.match(line) or REQUIRE.match(line)
            if match:
                modules[match.group(1)] = match.group(2)
    return modules


def read_notices():
    """Return the module rows of the notices as {module: release}.

    The release sits in the First version column. It reads vX.Y.Z for a
    release. A row for a dependency that no release holds yet is allowed
    and carries no release, so it maps to None.
    """
    rows = {}
    in_go_table = False
    with open(NOTICES, encoding="utf-8") as handle:
        for number, line in enumerate(handle, start=1):
            if line.startswith("## Go Dependencies"):
                in_go_table = True
                continue
            if in_go_table and line.startswith("## "):
                in_go_table = False
            if not in_go_table or not line.startswith("|"):
                continue
            match = ROW.match(line)
            if not match:
                continue
            module = match.group(1)
            cells = [cell.strip() for cell in line.strip().strip("|").split("|")]
            if set(cells[0]) <= {"-", " "} or cells[0].startswith("--"):
                continue
            if len(cells) < 4:
                print(f"{NOTICES}:{number}: row needs four cells: {module}")
                sys.exit(1)
            rows[module] = cells[3]
    return rows


def main():
    if not re.search(r"^go\s+1\.", open(GO_MOD, encoding="utf-8").read(), re.M):
        print(f"{GO_MOD}: no go directive found")
        return 1
    go_mod = read_go_mod()
    notices = read_notices()
    problems = 0
    # Each module in go.mod needs a row in the notices.
    for module in sorted(go_mod):
        if module not in notices:
            print(f"{GO_MOD}: {module} has no row in {NOTICES}")
            problems += 1
    # A row for a shipped module must name the release that first ships it.
    for module, release in sorted(notices.items()):
        if module not in go_mod:
            # The row names a dependency of a later release.
            continue
        if not RELEASE.match(release):
            print(
                f"{NOTICES}: {module} row must name the release, got {release!r}"
            )
            problems += 1
    print(f"checked {len(go_mod)} go.mod module(s) against "
          f"{len(notices)} notice row(s), {problems} problem(s)")
    return 1 if problems else 0


if __name__ == "__main__":
    sys.exit(main())
