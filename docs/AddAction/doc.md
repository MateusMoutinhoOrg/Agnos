# Add an Action

## Description
Covers adding a reusable operation to Agnos — the logic layer under `sandbox/internal/actions/<name>/` that a command calls and that `api.Sandbox.Actions` exposes to Go callers. Exposing it on the command line is [AddAgnosCommand](/docs/AddAgnosCommand/doc.md); the transactional filesystem every action writes through is [SmartIO](/docs/SmartIO/doc.md); the shape the two files must have is the Action specification in [Specs](/docs/Specs/doc.md).

### Rules
- An action is two files: `<name>.go`, holding the public entry that opens a SmartIO, calls the internal function, persists and runs any follow-up; and `<name>_internal.go`, holding the logic against an already-open `*smartio.SmartIO` so another action can compose it inside its own transaction.
- Every path handed to SmartIO is **project-relative** (`"go.mod"`, `"sandbox/api"`), never `path + "/..."`. SmartIO joins the root on at the boundary.
- Nothing hits disk before `io.Persist()`. An action that then runs `build` as a follow-up must persist **first**: `build`'s collectors list directories from disk and would not see pending writes.
- Progress goes through `deps.Std.Log`; failures come back as an `error` built with `deps.Std.Errorf`. An action never calls `Printf` — the result is the command's to print.
- A follow-up build names its runtime through `api.BuildProps`: actions that add something pass `api.RuntimeGo`, actions that remove something pass `api.RuntimeNone`.
- Register the action in `sandbox/api/actions.go` **and** `sandbox/binds/actions.go` in the same commit. `verify` requires every `binds/` file to mirror an `api/` file and declare only functions.

---

## Workflow
1. Create the package with its two files. Take the props an action needs as plain parameters, or as a props struct declared in `sandbox/api` when there are more than three — `api.StartProps`, `api.FieldProps` and `api.CommandProps` are the existing ones:
   ```bash
   mkdir sandbox/internal/actions/rename_command
   ```
2. Write the internal function against an open SmartIO. Read through `io`, mutate, write back with `io.WriteFile` / `io.WriteFileOverwrite` / `io.RemoveDir`; reuse the helpers in `sandbox/internal/utils` rather than reparsing YAML by hand:
   ```go
   // sandbox/internal/actions/rename_command/rename_command_internal.go
   package rename_command

   import (
       "strings"

       "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
   )

   // RenameCommandInternal moves sandbox/internal/commands/<from>/ to <to>/
   // inside the open transaction; the caller persists and rebuilds.
   func RenameCommandInternal(deps *deps.Deps, io *smartio.SmartIO, path string, from string, to string) error {
       deps.Std.Log("rename-command %s -> %s \n", from, to)

       fromDir := utils.CommandDir(from)
       toDir := utils.CommandDir(to)
       if io.IsDir(toDir) {
           return deps.Std.Errorf("command %q already exists", to)
       }
       for _, file := range io.ListFiles(fromDir) {
           content, err := io.ReadFile(file)
           if err != nil {
               return err
           }
           base := file[strings.LastIndex(file, "/")+1:]
           if err := io.WriteFile(toDir+"/"+base, content); err != nil {
               return err
           }
       }
       io.RemoveDir(fromDir)
       return nil
   }
   ```
3. Write the public entry: open the transaction rooted at the project, call the internal function, persist, then run the follow-up build with the runtime the change deserves:
   ```go
   // sandbox/internal/actions/rename_command/rename_command.go
   package rename_command

   import (
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
       buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
   )

   func RenameCommand(deps *deps.Deps, path string, from string, to string) error {
       io := smartio.New(deps, path, config.ProjectName)
       if err := RenameCommandInternal(deps, io, path, from, to); err != nil {
           return err
       }
       if err := io.Persist(); err != nil {
           return err
       }
       return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
   }
   ```
4. Declare the field on the contract, in `sandbox/api/actions.go`, next to the others — `api` imports nothing, so the signature uses only builtin types and other `api` structs:
   ```go
   RenameCommand func(path string, from string, to string) error
   ```
5. Bind it in `sandbox/binds/actions.go`, the one place every action is assigned:
   ```go
   sandbox.Actions.RenameCommand = func(path string, from string, to string) error {
       return renameCommandAction.RenameCommand(deps, path, from, to)
   }
   ```
6. Compile and verify the schema still holds:
   ```bash
   go build ./cmd/... ./sandbox/... ./adapters/... && go run ./cmd/main verify
   ```
7. Expose it on the command line following [AddAgnosCommand](/docs/AddAgnosCommand/doc.md), and list it in the **Actions** section of [PublicApi](/docs/PublicApi/doc.md) with a detail page — required by [HandleDocuments](/docs/HandleDocuments/doc.md).
