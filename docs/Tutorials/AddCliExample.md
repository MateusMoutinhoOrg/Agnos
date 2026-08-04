# Add a CLI Example

## Description
Covers creating a shell script in [cliExamples/](/cliExamples/) that demonstrates one goal against the built CLI. The Go counterpart — a program wiring the library from code — is covered by [AddSample.md](/docs/Tutorials/AddSample.md).

### Rules
- The script must follow its specification — locate it in [Specs.md](/docs/References/Specs.md).
- It must run from the project root with no arguments and no prior setup, and must never write outside its own scratch directory.
- Adding one requires updating the CLI Examples section of [README.md](/README.md) and [Structure.md](/docs/References/Structure.md).

---

## Workflow
1. Create the script under [cliExamples/](/cliExamples/), named `example<N>.sh` with the next free number — the numbering is the reading order, so a script demonstrating a later idea gets a later number.
2. Open it with the shebang, a comment naming the goal, how to run it, and the shell options:
   ```bash
   #!/usr/bin/env bash
   # example4.sh — one sentence naming what this script demonstrates.
   #
   # Run it from the project root:
   #   bash ./cliExamples/example4.sh
   set -euo pipefail
   ```
3. Build the CLI into a scratch directory and point it at a budget of its own, so nothing the script does touches the user's records:
   ```bash
   workdir="$(mktemp -d)"
   trap 'rm -rf "$workdir"' EXIT

   go build -o "$workdir/agnos" ./cmd/main
   export AGNOS_DATA="$workdir/data"
   agnos() { "$workdir/agnos" "$@"; }
   ```
4. Write the demonstration as sections, each announced by an `echo` line saying what the commands below it show:
   ```bash
   echo "== record what came in and what went out"
   agnos --quiet category add groceries
   agnos spend groceries "weekly shopping" 84.50
   ```
5. Make it executable and run it, following [RunCliExample.md](/docs/Tutorials/RunCliExample.md):
   ```bash
   chmod +x ./cliExamples/example4.sh
   bash ./cliExamples/example4.sh
   ```
6. Add the script to the CLI Examples section of [README.md](/README.md).
7. Register it in [Structure.md](/docs/References/Structure.md) if it introduces anything the schema does not already describe.
