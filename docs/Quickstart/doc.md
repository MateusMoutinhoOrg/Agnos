# Quickstart

Every command takes the project dir via `--path` (default `.`). Every step ends by running `build`, so the tree always compiles.

```bash
agnos start --project-name my-tool --module github.com/you/my-tool   # AgnosConfig/, go.mod, sandbox skeleton
agnos list-extensions                                                 # what agnos generates for this project
agnos deps-init                                                       # sandbox/deps/ + adapters/ (sandbox-deps: true)
agnos add-dep iodeps                                              # any name from `agnos list-deps`
agnos cli-init                                                        # cmd/main, dispatch, help, version (sandbox-cli: true)
agnos add-command greet --help "Say hello" --category Demo
agnos add-flag name --command greet --identifier --name --identifier -n --required --description "who to greet"
agnos add-arg times --command greet --type int --min 1 --default 1 --description "how many times"
```

Write the one hand-written file, `sandbox/internal/commands/greet/handler.go`:

```go
package greet

import (
	"github.com/you/my-tool/sandbox/api"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	name := command.GetString("name")
	sandbox.Deps.Std.Log("greeting %s\n", name)        // stderr, silenced by --quiet
	for i := 0; i < command.GetInt("times"); i++ {
		sandbox.Deps.Std.Printf("hello, %s\n", name)   // stdout, the result
	}
	return api.ExitOk
}
```

What a handler may and may not do is in [Rules](../Rules/doc.md#handlers).

Then:

```bash
agnos build                         # verify + regenerate + go build
go run ./cmd/main greet -n bob 2
agnos compile --target all          # release/<target> per platform
agnos publish --draft               # build, compile all, `gh release create <version>`
```

What agnos generates is yours to choose — `agnos disable-extension readme` hands `README.md`
back to you, and the next `build` leaves it alone. Every key is in
[Extensions](../Extensions/doc.md).
