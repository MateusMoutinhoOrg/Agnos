# Shape a Command from the Command Line

## Description
Covers declaring the whole surface of a command — its flags, positional arguments, aliases, help text and visibility — through `agnos` commands, so `entries.yaml` is never edited by hand. Creating the command in the first place is step 5 of [ScaffoldProject.md](/docs/Tutorials/ScaffoldProject.md); what every key means, and how it becomes a Go field, is [EntriesYaml.md](/docs/References/EntriesYaml.md); writing the handler that receives the values is [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md).

### Rules
- `--command` names the command receiving the field, by identifier or by package name (`greet` or `sandbox/internal/commands/greet`'s `greet`); the project is taken from `--path`, defaulting to the current directory.
- Every one of these commands rewrites `entries.yaml` from its parsed form and then runs `build`, so the generated `entries.go`, the dispatch and `help` never lag behind the declaration. Because the file is re-rendered, YAML comments are dropped and keys come out in alphabetical order — the content round-trips, only the layout is canonicalized.
- `--required` is refused on a boolean flag, and on any field that has a `--default`: absence is already covered in both cases.
- `min` and `max` apply to `int` and `float` fields only; the generated dispatch rejects an out-of-range value with a usage error before the handler runs.
- An `array` argument collects every remaining positional, so it must stay the **last** argument.
- The `help` command belongs to `agnos`: `add-command help`, `remove-command help`, and hand edits to its `handler.go` are refused or overwritten — see [CommandDispatch.md](/docs/References/CommandDispatch.md#the-help-command).

---

## Workflow

### Add a flag

1. Add a flag with one or more spellings. Without `--identifier` the flag answers to `--<name>`; the name is also what the generated `Entries` field is derived from (`out-file` becomes `OutFile`):
   ```bash
   agnos add-flag output --command exec --identifier --out --identifier -o --type string --required --description "where the output is written"
   ```
2. Add a boolean flag. Booleans are presence flags: no value follows them, and they cannot be `--required`:
   ```bash
   agnos add-flag verbose --command exec --type boolean --description "print every step"
   ```
3. Add a numeric flag with a default and a range. `--default`, `--min` and `--max` are typed as raw literals so "unset" stays distinguishable from zero:
   ```bash
   agnos add-flag retries --command exec --type int --min 0 --max 5 --default 1 --description "attempts before giving up"
   ```
4. Add a repeatable flag. `--array` turns the field into a slice collecting every occurrence:
   ```bash
   agnos add-flag tag --command exec --type string --array --example "--tag a --tag b"
   ```
5. Declare a boolean flag named `quiet` to get a working `--quiet`. The generated dispatch silences the progress channel as soon as it has read that flag; the handler has nothing to do:
   ```bash
   agnos add-flag quiet --command exec --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
   ```

### Add a positional argument

6. Append an argument. Positional arguments bind by declaration order, so the first `add-arg` is the first word after the command:
   ```bash
   agnos add-arg file --command exec --type string --required --description "the file to process"
   ```
7. Insert an argument at a given index with `--position` (zero-based; the default `-1` appends). Later arguments shift down:
   ```bash
   agnos add-arg count --command exec --type int --min 1 --default 1 --position 0
   ```

### Rewrite the command-level keys

8. Change the help line, the category `help` groups it under, or the long description shown by `help <command>`. Keys not passed are left untouched:
   ```bash
   agnos set-command exec --help "run the thing" --category Core --long-description "Runs the thing end to end."
   ```
9. Add an alias or a usage example. `--identifier` and `--example` append and deduplicate, so running them twice is harmless:
   ```bash
   agnos set-command exec --identifier run --example "exec file.txt --out result.txt"
   ```
10. Hide the command from the `help` listings, or show it again. A hidden command still dispatches and still answers `help <command>`:
    ```bash
    agnos set-command exec --hidden
    agnos set-command exec --visible
    ```

### Remove a field or a command

11. Drop a flag by its name or by any of its identifiers, or drop an argument by its name. The follow-up build runs with the `none` runtime, because hand-written handler code may still refer to the removed field — fix `handler.go`, then run `agnos build` yourself:
    ```bash
    agnos remove-flag --out --command exec
    agnos remove-arg count --command exec
    agnos build
    ```
12. Delete a whole command. The package directory goes — `entries.yaml`, `entries.go`, `handler.go` and anything else inside — and the dispatch and `help` forget it on the rebuild:
    ```bash
    agnos remove-command exec
    ```
