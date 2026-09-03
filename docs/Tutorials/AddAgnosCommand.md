# Add a Command to Agnos

## Description
Covers exposing an action on `agnos`'s own command line. Agnos-Cli is itself a project of the shape it generates, so a command here is the same three files as everywhere else — `entries.yaml`, generated `entries.go`, hand-written `handler.go` under `sandbox/internal/commands/<name>/` — and is declared with the same commands a user runs, pointed at this checkout. The logic the command calls is [AddAction.md](/docs/Tutorials/AddAction.md); the generic handler rules are [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md).

### Rules
- Run the **bootstrap binary**, never an installed `agnos`, against this checkout — see [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md).
- Every command that operates on a project takes the directory through a `--path` flag with `default: "."`, never a positional, so `entries.Path` is always populated.
- Every command declares a boolean `quiet` flag with identifiers `--quiet` and `-q`; the dispatch does the silencing.
- The handler calls the action package directly and maps its outcome to an exit code: `api.ExitFailure` on an error, `api.ExitOk` otherwise. It prints the result, if the action returns one, through `deps.Std.Printf`.
- Adding a command requires updating [Commands.md](/docs/References/Commands.md) in the same commit; `help` regenerates itself.

---

## Workflow
1. Compile the bootstrap binary from the current tree:
   ```bash
   go build -o release/bootstrap.bin ./cmd/main
   ```
2. Declare the command with it. `--help` and `--category` are required; the categories in use are `Core Commands`, `Cli System`, `Dependencies`, `Dependency System` and `Info`:
   ```bash
   ./release/bootstrap.bin add-command rename-command --help "Rename a command package" --category "Cli System"
   ```
3. Grow the declaration from the command line — the `--path` and `quiet` flags every command carries, then the command's own fields:
   ```bash
   B=./release/bootstrap.bin
   $B add-flag path --command rename-command --type string --default . --description "the dir holding the project (defaults to the current directory)"
   $B add-flag quiet --command rename-command --identifier --quiet --identifier -q --type boolean --description "Quiets the cli output"
   $B add-arg from --command rename-command --type string --required --description "the command to rename"
   $B add-arg to --command rename-command --type string --required --description "its new name"
   $B set-command rename-command --long-description "Moves sandbox/internal/commands/<from>/ to <to>/ and runs build." --example "rename-command exec run"
   ```
4. Replace the stub in `sandbox/internal/commands/rename_command/handler.go` with the call into the action:
   ```go
   package rename_command

   import (
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
       "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
       renameCommandAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/rename_command"
   )

   func CommandHandler(deps *deps.Deps, entries *Entries) int {
       if err := renameCommandAction.RenameCommand(deps, entries.Path, entries.From, entries.To); err != nil {
           deps.Std.Error("%v\n", err)
           return api.ExitFailure
       }
       return api.ExitOk
   }
   ```
5. Regenerate with the bootstrap binary — `entries.go`, `climain.go` and `help` catch up — and try the command:
   ```bash
   $B build
   go run ./cmd/main help rename-command
   ```
6. Add the command's row to the table of [Commands.md](/docs/References/Commands.md), in its category, with every flag and argument.
7. Finish with the bootstrap loop of [BootstrapAgnos.md](/docs/Tutorials/BootstrapAgnos.md) so the installed `agnos` carries the new command.
