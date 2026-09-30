# Quickstart

Every command takes the project dir via `--path` (default `.`). Every step ends by running `build`, so the tree always compiles.
Typing none of it: `agnos interview` asks these same steps as questions, one suggested at a time — [Interview](../Interview/doc.md).

```bash
agnos start --project-name my-tool --module github.com/you/my-tool   # AgnosConfig/, go.mod, sandbox skeleton
agnos list-extensions                                                 # what agnos generates for this project
agnos deps-init                                                       # sandbox/deps/ + adapters/ (sandbox-deps: true)
agnos add-dep iodeps                                              # any name from `agnos list-deps`
agnos cli-init                                                        # cmd/main, dispatch, help, version (sandbox-cli: true)
agnos add-command greet --help "Say hello" --category Demo
agnos add-flag name --command greet --key --name --key -n --required --description "who to greet"
agnos add-arg times --command greet --type integer --default 1 --description "how many times"
```

Write the one hand-written file, `sandbox/internal/commands/greet/InternalPureHandler.go`
(`add-command` wrote a stub; `build` generated `entries.go` beside it):

```go
package greet

import (
	"github.com/you/my-tool/sandbox/api"
	"github.com/you/my-tool/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	response.Log("greeting %s\n", entries.Name)         // stderr
	for i := 0; i < entries.Times; i++ {
		response.Printf("hello, %s\n", entries.Name)   // stdout, the result
	}
	return nil
}
```

What a handler may and may not do is in [Rules](../Rules/doc.md#handlers).

Then:

```bash
agnos build                         # verify + regenerate + go build
go run ./cmd/main greet 2 -n bob    # args first: tokens after the flags are not read as args
agnos compile --target all          # release/<target> per platform
agnos publish --draft --release-name v0.1.0   # build, compile all, `gh release create v0.1.0`
```

`publish` names the release after `version` in `AgnosConfig/project.yaml` when no
`--release-name` is given; `start` writes it as `null`, which `publish` refuses.

What agnos generates is yours to choose — `agnos disable-extension readme` hands `README.md`
back to you, and the next `build` leaves it alone. Every key is in
[Extensions](../Extensions/doc.md).
