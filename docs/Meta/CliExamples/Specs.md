# CliExamples Specification

## Description
Defines the required shape of a runnable command-line example in `cliExamples/<Name>.sh`. A CLI example is a self-contained shell script that builds the binary in [cmd/main](/cmd/main/) and drives it the way a user would from a terminal — one script per goal a person actually has. Its Go counterpart, which wires the library from code instead, is governed by [LibraryExamples](/docs/Meta/LibraryExamples/Specs.md).

### Rules
- Each CLI example is a single `.sh` file directly under `cliExamples/`, named with a descriptive PascalCase name matching the goal it demonstrates — e.g. `ManageCategories.sh`, `TrackTransactions.sh`.
- The script opens with a shebang (`#!/usr/bin/env bash`), a comment block naming the goal it demonstrates, and the command that runs it from the project root.
- It must be runnable with no arguments and no prior setup, from the project root: `bash ./cliExamples/<Name>.sh`.
- It builds the CLI itself with `go build -o "$workdir/agnos" ./cmd/main` into a `mktemp -d` directory removed by a `trap` on exit, and points the binary at a budget of its own by exporting `AGNOS_DATA`. A CLI example never writes to the records in the user's home directory and never requires the CLI to be installed first.
- Every command it runs goes through the built binary — the example demonstrates the interface, so it must never `go run` a library example or reach into the sandbox.
- Sections are announced with an `echo "== …"` line saying what the commands below show, so the transcript reads on its own.
- Adding, renaming, or deleting a CLI example requires updating the CLI Examples section of [README.md](/README.md) and [Structure.md](/docs/Structure.md) — see [HandleCliExamples.md](/docs/HandleCliExamples.md).

## Structure
1. **Shebang**: `#!/usr/bin/env bash`.
2. **Header comment**: the script's name, the goal it demonstrates, and how to run it.
3. **Shell options**: `set -euo pipefail`, or `set -uo pipefail` when the example deliberately inspects a failing exit code.
4. **Scratch setup**: a `mktemp -d` work directory with a cleanup `trap`, the `go build` of `./cmd/main` into it, an exported `AGNOS_DATA`, and an `agnos()` shell function calling the built binary.
5. **Demonstration sections**: each an `echo "== …"` line followed by the commands it describes.

> **Note**: For a concrete example, refer to [sample.sh](./sample.sh).
