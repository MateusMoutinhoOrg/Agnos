# Handle CLI Examples

## Description
Covers creating and running shell scripts in [examples/cliExamples/](/examples/cliExamples/) that demonstrate how a user would drive the built CLI from a terminal — typically after adding a command through [HandleCliCommands.md](/docs/HandleCliCommands.md). The Go counterparts, wiring the library from code, are covered by [HandleSamples.md](/docs/HandleSamples.md).

---

## Run a CLI Example

### Rules
- A CLI example needs nothing installed: it builds the binary itself.
- It writes to a scratch directory it removes on exit, so running one never touches the records in your home directory.

### Workflow
1. Browse [examples/cliExamples/](/examples/cliExamples/) and pick a script — each is named after the goal it demonstrates, so `ManageCategories.sh` is a good starting point.
2. Run it from the project root:
   ```bash
   bash ./examples/cliExamples/ManageCategories.sh
   ```
3. Read the transcript alongside the script: each `== …` line in the output is the comment above the commands that produced what follows it.
4. Run the rest in order to see the whole interface:
   ```bash
   for script in ./examples/cliExamples/*.sh; do bash "$script"; done
   ```
5. Try the same commands against your own budget once you have installed the binary, following [InstallCli.md](/docs/InstallCli.md) and [UseCli.md](/docs/UseCli.md).

---

## Add a CLI Example

### Rules
- The script must follow its specification — locate it in [Specs.md](/docs/Specs.md).
- It must run from the project root with no arguments and no prior setup, and must never write outside its own scratch directory.
- Adding one requires updating the CLI Examples section of [README.md](/README.md) and [Structure.md](/docs/Structure.md).

### Workflow
1. Create the script under [examples/cliExamples/](/examples/cliExamples/), named with a descriptive PascalCase name matching the goal it demonstrates — e.g. `ManageCategories.sh`, `TrackTransactions.sh`.
2. Open it with the shebang, a comment naming the goal, how to run it, and the shell options:
   ```bash
   #!/usr/bin/env bash
   # <Name>.sh — one sentence naming what this script demonstrates.
   #
   # Run it from the project root:
   #   bash ./examples/cliExamples/<Name>.sh
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
5. Make it executable and run it:
   ```bash
   chmod +x ./examples/cliExamples/<Name>.sh
   bash ./examples/cliExamples/<Name>.sh
   ```
6. Add the script to the CLI Examples section of [README.md](/README.md).
7. Register it in [Structure.md](/docs/Structure.md) if it introduces anything the schema does not already describe.
