# Library Initialization

## Description
Covers installing Agnos as a Go module and running its actions from a program of your own instead of from the terminal: the same `start`, `build` and `dep-install` the CLI exposes are fields of `api.Sandbox.Actions`. Shaping a project's commands from Go is [ShapeProjectFromGo.md](/docs/Tutorials/ShapeProjectFromGo.md); every action is listed in [PublicApi.md](/docs/References/PublicApi.md#actions); building the deps differently is [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md).

### Rules
- The library is wired the same way the binary is: `standard.New()` builds a `deps.Deps`, `sandbox.New(&deps)` injects it, and the returned `*api.Sandbox` carries every behavior — `Actions` for the operations, `Cli` for the whole command-line interface.
- Every action takes the target project directory as a parameter (`Path`) and scopes every read and write to it, exactly as the `--path` flag does.
- Actions report progress through `deps.Std.Log` (standard error) and return an `error` on failure; they print no result of their own except `DepList`, which returns it.

---

## Workflow
1. Install the module:
   ```bash
   go get github.com/MateusMoutinhoOrg/Agnos@latest
   ```
2. Create a `main.go` that scaffolds a project the way `agnos start` does. `Start` writes the configuration, renders the skeleton, persists, and hands the result to the Go toolchain, so the directory it leaves behind compiles:
   ```go
   package main

   // 1. Import the standard adapter, the sandbox and its contract package
   import (
       "log"

       agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
       agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
       "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
   )

   func main() {
       // 2. Build the deps through an adapter, and inject them
       deps := agnosadapter.New()
       lib := agnoslib.New(&deps)

       // 3. Run an action — the same code `agnos start` runs
       module := "github.com/you/my-tool"
       err := lib.Actions.Start(api.StartProps{
           Path:        "./my-tool",
           ProjectName: "my-tool",
           Module:      &module,
       })
       if err != nil {
           log.Fatal(err)
       }
   }
   ```
3. Run it, then look at what it wrote:
   ```bash
   mkdir my-tool && go run main.go && ls my-tool
   ```
4. Chain the other actions the same way. `DepsInit`, `DepInstall` and `CliInit` each end by running `build` themselves, exactly like their commands:
   ```go
   if err := lib.Actions.DepsInit("./my-tool"); err != nil {
       log.Fatal(err)
   }
   if err := lib.Actions.DepInstall("./my-tool", "iodeps"); err != nil {
       log.Fatal(err)
   }
   if err := lib.Actions.CliInit("./my-tool"); err != nil {
       log.Fatal(err)
   }
   ```
5. Rebuild or check a project by hand. `Build` takes `api.BuildProps`, whose `Runtime` is `api.RuntimeGo` (tidy and compile) or `api.RuntimeNone` (render only); `Verify` takes the path alone and writes nothing:
   ```go
   if err := lib.Actions.Verify("./my-tool"); err != nil {
       log.Fatal(err)
   }
   if err := lib.Actions.Build(api.BuildProps{Path: "./my-tool", Runtime: api.RuntimeGo}); err != nil {
       log.Fatal(err)
   }
   ```
6. Run the whole command-line interface over an argument vector when you want the CLI's behavior — its parsing, its usage errors, its exit code — without a process:
   ```go
   code := lib.Cli.CliMain([]string{"dep-list", "--path", "./my-tool"})
   log.Println("exit", code)
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

	path := "./my-tool"
	module := "github.com/you/my-tool"

	steps := []func() error{
		func() error {
			return lib.Actions.Start(api.StartProps{Path: path, ProjectName: "my-tool", Module: &module})
		},
		func() error { return lib.Actions.DepsInit(path) },
		func() error { return lib.Actions.DepInstall(path, "iodeps") },
		func() error { return lib.Actions.CliInit(path) },
		func() error { return lib.Actions.Verify(path) },
		func() error { return lib.Actions.Build(api.BuildProps{Path: path, Runtime: api.RuntimeGo}) },
	}
	for _, step := range steps {
		if err := step(); err != nil {
			log.Fatal(err)
		}
	}

	// The same interface the binary exposes, as a function call.
	code := lib.Cli.CliMain([]string{"dep-list", "--path", path})
	log.Println("dep-list exited with", code)
}
```
