# LibUsage

```bash
go get github.com/MateusMoutinhoOrg/Agnos@latest
```

```go
import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

deps := standard.New()          // every adapter lib bound
lib := sandbox.New(&deps)       // *api.Sandbox: Actions + Cli

module := "github.com/you/my-tool"
err := lib.Actions.Start(api.StartProps{Path: "./my-tool", ProjectName: "my-tool", Module: &module})
err = lib.Actions.DepsInit("./my-tool")
err = lib.Actions.DepInstall("./my-tool", "iodeps")
err = lib.Actions.CliInit("./my-tool")
err = lib.Actions.AddCommand("./my-tool", "greet", "Say hello", "Demo")
err = lib.Actions.AddFlag(api.FieldProps{Path: "./my-tool", Command: "greet", Name: "name", Identifiers: []string{"--name", "-n"}, Required: true})
err = lib.Actions.Build(api.BuildProps{Path: "./my-tool", Runtime: api.RuntimeGo})

code := lib.Cli.CliMain([]string{"dep-list", "--path", "./my-tool"})   // the whole CLI, in-process
```

- Actions mirror the CLI one to one; signatures in [PublicApi](/docs/PublicApi/doc.md#actions). Each takes the project dir, scopes every read/write to it, logs progress via `deps.Std.Log`, returns `error`.
- Every action that adds something ends by running `build`, like its command.

## Custom Deps

Patch fields **before** `sandbox.New(&deps)`; the binders capture the pointer. Start from `standard.New()`: an unfilled field is a nil func that panics on first call.

```go
deps := standard.New()
var out bytes.Buffer
deps.Std.Printf = func(f string, a ...any) (int, error) { return fmt.Fprintf(&out, f, a...) }
lib := sandbox.New(&deps)
```

A permanent mix is its own `adapters/availables/<name>/new.go` binding only the libs you want (`iodeps.Bind(&deps)`, ...). `standard/new.go` is regenerated on every build; other dirs under `availables/` are left alone. Wire it from a `cmd/<name>/main.go` of your own, since `cmd/main/main.go` is generated.
