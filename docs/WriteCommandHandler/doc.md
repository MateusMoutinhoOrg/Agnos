# Write a Command Handler

## Description
Covers filling in `sandbox/internal/commands/<name>/handler.go` — the one file of a command you write by hand — in a project scaffolded by `agnos`. The declaration next to it is grown from the command line by [ShapeCommands](/docs/ShapeCommands/doc.md); how the generated dispatch turns that declaration into the `Entries` the handler receives is [CommandDispatch](/docs/CommandDispatch/doc.md); an effect the handler needs and no installed dep provides is [AddAdapterLib](/docs/AddAdapterLib/doc.md).

### Rules
- The handler has exactly one exported function, `CommandHandler(deps *deps.Deps, entries *Entries) int`. `agnos build` generates a call to it under that name and no other.
- The handler lives inside the closed sandbox: it may import `sandbox/api`, `sandbox/deps` and anything under `sandbox/internal`, and nothing outside `sandbox/` — no `os`, no `fmt`, no third-party module. `verify` rejects the import; see [SandboxIsolation](/docs/SandboxIsolation/doc.md).
- Every effect goes through a `deps` field named after its `sandbox/deps/<x>/` directory: `deps.Std`, `deps.Iodeps`, `deps.Requestdeps`, … — see [PublicApi](/docs/PublicApi/doc.md#dependency-contracts).
- Three output channels, one meaning each: `deps.Std.Printf` writes the command's **result** to standard output, `deps.Std.Log` writes **progress** to standard error and is what `--quiet` silences, `deps.Std.Error` writes **failures** to standard error.
- Return `api.ExitOk` (`0`) on success and `api.ExitFailure` (`1`) when a well-formed command could not be carried out. Never return `api.ExitUsage` (`2`): a wrong command line is rejected by the generated dispatch before the handler runs, so a handler that starts is holding valid input.
- `entries.go` is regenerated on every build — put nothing in it.

---

## Workflow
1. Open the stub `add-command` wrote. It compiles and prints one line; everything below replaces its body:
   ```go
   package greet

   import (
       "github.com/you/my-tool/sandbox/api"
       "github.com/you/my-tool/sandbox/deps"
   )

   func CommandHandler(deps *deps.Deps, entries *Entries) int {
       deps.Std.Printf("greet called\n")
       return api.ExitOk
   }
   ```
2. Read the fields `agnos build` generated into `entries.go` — flags first, then positional arguments, each in declaration order, each already typed and range-checked. A field with a `default` is always populated; an `array` field is a slice:
   ```go
   // sandbox/internal/commands/greet/entries.go — generated, do not edit
   type Entries struct {
       Name  string // --name / -n, required
       Times int    // positional, default 1, min 1
   }
   ```
3. Write the logic against those fields. Report progress through `Log`, so `--quiet` can turn it off, and the result through `Printf`, so a script reading standard output gets only the result:
   ```go
   func CommandHandler(deps *deps.Deps, entries *Entries) int {
       deps.Std.Log("greeting %s %d time(s)\n", entries.Name, entries.Times)
       for i := 0; i < entries.Times; i++ {
           deps.Std.Printf("hello, %s\n", entries.Name)
       }
       return api.ExitOk
   }
   ```
4. Reach the outside world through `deps`, never through the standard library. A handler that writes a file asks `deps.Iodeps` — installed by `agnos dep-install iodeps` — and reports a failure with `Error` and `ExitFailure`:
   ```go
   func CommandHandler(deps *deps.Deps, entries *Entries) int {
       content := []byte("hello, " + entries.Name + "\n")
       if err := deps.Iodeps.WriteFile(entries.Output, content); err != nil {
           deps.Std.Error("could not write %s: %v\n", entries.Output, err)
           return api.ExitFailure
       }
       deps.Std.Log("wrote %s\n", entries.Output)
       return api.ExitOk
   }
   ```
5. Keep logic worth reusing out of the handler. A package under `sandbox/internal/<name>/` is reachable from every handler and from nothing outside the sandbox; the handler then parses nothing and only calls it and reports.
6. Rebuild and run. `build` compiles the whole project, so a handler that does not compile fails here, not at first use:
   ```bash
   agnos build
   go run ./cmd/main greet -n bob 2
   ```
7. Try the rejections the dispatch makes for you — none of them reaches the handler, and each exits `2`:
   ```bash
   go run ./cmd/main greet              # required flag 'name' not provided
   go run ./cmd/main greet -n bob 0     # arg 'times' must be >= 1
   go run ./cmd/main greet -n bob x     # arg 'times': "x" is not a valid integer
   go run ./cmd/main greet -n bob 1 x   # unexpected argument "x"
   ```

## Full Code

```go
package greet

import (
	"github.com/you/my-tool/sandbox/api"
	"github.com/you/my-tool/sandbox/deps"
)

// CommandHandler backs the `greet` verb. Entries.Name comes from --name / -n
// and is required; Entries.Times is the optional positional, at least 1,
// defaulting to 1. Both are validated before this function runs.
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deps.Std.Log("greeting %s %d time(s)\n", entries.Name, entries.Times)
	for i := 0; i < entries.Times; i++ {
		deps.Std.Printf("hello, %s\n", entries.Name)
	}
	return api.ExitOk
}
```
