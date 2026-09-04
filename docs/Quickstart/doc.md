# Quickstart

Every command takes the project dir via `--path` (default `.`). Every step ends by running `build`, so the tree always compiles.

```bash
agnos start --project-name my-tool --module github.com/you/my-tool   # AgnosConfig/, go.mod, sandbox skeleton
agnos deps-init                                                       # sandbox/deps/ + adapters/
agnos dep-install iodeps                                              # any name from `agnos dep-list`
agnos cli-init                                                        # cmd/main, dispatch, help, version (installs std + argvdeps)
agnos add-command greet --help "Say hello" --category Demo
agnos add-flag name --command greet --identifier --name --identifier -n --required --description "who to greet"
agnos add-arg times --command greet --type int --min 1 --default 1 --description "how many times"
```

Write the one hand-written file, `sandbox/internal/commands/greet/handler.go`:

```go
package greet

import (
	"github.com/you/my-tool/sandbox/api"
	"github.com/you/my-tool/sandbox/deps"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deps.Std.Log("greeting %s\n", entries.Name)          // stderr, silenced by --quiet
	for i := 0; i < entries.Times; i++ {
		deps.Std.Printf("hello, %s\n", entries.Name)     // stdout, the result
	}
	return api.ExitOk
}
```

Rules for a handler:
- Only `CommandHandler(deps *deps.Deps, entries *Entries) int` is exported. `Entries` is generated from `entries.yaml` (flags first, then args, declaration order), already typed, defaulted and range-checked.
- Import nothing outside `sandbox/`. Every effect goes through `deps.<Contract>` (`deps.Iodeps`, `deps.Std`, ...), see [PublicApi](../PublicApi/doc.md).
- Return `api.ExitOk` or `api.ExitFailure`. Never `api.ExitUsage`: the dispatch rejects bad input before the handler runs.
- Reusable logic goes in `sandbox/internal/<pkg>/`.

Then:

```bash
agnos build                         # verify + regenerate + go build
go run ./cmd/main greet -n bob 2
agnos compile --target all          # release/<target> per platform
agnos publish --draft               # build, compile all, `gh release create <version>`
```
