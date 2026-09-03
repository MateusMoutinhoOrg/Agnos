# Shape a Project's Commands from Go

## Description
Covers declaring a scaffolded project's command surface from a Go program — the `AddCommand`, `AddFlag`, `AddArg`, `SetCommand` and `Remove*` actions behind the commands of [ShapeCommands.md](/docs/Tutorials/ShapeCommands.md). Initializing the library is [LibInitialization.md](/docs/Tutorials/LibInitialization.md); the props structs are detailed in [api.FieldProps](/docs/References/PublicApi/api.FieldProps.md) and [api.CommandProps](/docs/References/PublicApi/api.CommandProps.md).

### Rules
- The project must already have its CLI layer (`CliInit` has run): every editor loads `sandbox/internal/commands/<cmd>/entries.yaml` and fails when it is missing.
- `Default`, `Min` and `Max` on `api.FieldProps` are **strings** — the literal as it would be typed on the command line — so an unset value (`""`) stays distinguishable from `"0"`. `Position` is the index to insert at, and `-1` appends.
- Every editor runs `build` after writing, so the generated `entries.go`, dispatch and `help` are current when the action returns. Adding actions use the `go` runtime; removing actions render only, because a handler may still refer to the removed field.

---

## Workflow
1. Declare a command. `help` and `category` are required, as they are on the command line, and the name is normalized the same way (`My_Feature` becomes `my-feature`):
   ```go
   if err := lib.Actions.AddCommand(path, "greet", "Say hello", "Demo"); err != nil {
       log.Fatal(err)
   }
   ```
2. Add a flag through `api.FieldProps`. `Identifiers` default to `--<Name>` when left empty; `Type` defaults to `string`:
   ```go
   err := lib.Actions.AddFlag(api.FieldProps{
       Path:        path,
       Command:     "greet",
       Name:        "name",
       Identifiers: []string{"--name", "-n"},
       Type:        "string",
       Description: "who to greet",
       Required:    true,
       Position:    -1,
   })
   ```
3. Add a positional argument with a default and a lower bound. Both arrive as raw literals:
   ```go
   err = lib.Actions.AddArg(api.FieldProps{
       Path:        path,
       Command:     "greet",
       Name:        "times",
       Type:        "int",
       Description: "how many times",
       Default:     "1",
       Min:         "1",
       Position:    -1,
   })
   ```
4. Rewrite the command-level keys through `api.CommandProps`. Empty strings leave the current value alone; `Identifiers` and `Examples` append and deduplicate; `Hidden` and `Visible` are the two directions of one switch:
   ```go
   err = lib.Actions.SetCommand(api.CommandProps{
       Path:            path,
       Command:         "greet",
       LongDescription: "Greets someone, as many times as asked.",
       Identifiers:     []string{"hello"},
       Examples:        []string{"greet -n bob 2"},
   })
   ```
5. Remove a field or the whole command. A flag is matched by its name or by any of its identifiers:
   ```go
   err = lib.Actions.RemoveFlag(path, "greet", "-n")
   err = lib.Actions.RemoveArg(path, "greet", "times")
   err = lib.Actions.RemoveCommand(path, "greet")
   ```
6. Rebuild with the `go` runtime after a removal, once the handler no longer refers to what is gone:
   ```go
   err = lib.Actions.Build(api.BuildProps{Path: path, Runtime: api.RuntimeGo})
   ```

## Full Code

```go
package main

import (
	"log"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func main() {
	deps := agnosadapter.New()
	lib := agnoslib.New(&deps)
	path := "./my-tool" // already scaffolded, with its CLI layer installed

	must := func(err error) {
		if err != nil {
			log.Fatal(err)
		}
	}

	must(lib.Actions.AddCommand(path, "greet", "Say hello", "Demo"))

	must(lib.Actions.AddFlag(api.FieldProps{
		Path: path, Command: "greet", Name: "name",
		Identifiers: []string{"--name", "-n"}, Type: "string",
		Description: "who to greet", Required: true, Position: -1,
	}))

	must(lib.Actions.AddArg(api.FieldProps{
		Path: path, Command: "greet", Name: "times",
		Type: "int", Description: "how many times",
		Default: "1", Min: "1", Position: -1,
	}))

	must(lib.Actions.SetCommand(api.CommandProps{
		Path: path, Command: "greet",
		LongDescription: "Greets someone, as many times as asked.",
		Identifiers:     []string{"hello"},
		Examples:        []string{"greet -n bob 2"},
	}))
}
```
