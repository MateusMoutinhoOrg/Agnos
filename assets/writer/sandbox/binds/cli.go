package binds

import (
	api "{{.Module}}/sandbox/api"
	deps "{{.Module}}/sandbox/deps"
	cli "{{.Module}}/sandbox/internal/cli"
)

// CliBind wires the generated cli.CliMain (built by `agnos build` from every
// command's entries.yaml) onto the sandbox.
func CliBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Cli.CliMain = func(args []string) int {
		return cli.CliMain(deps, args)
	}
}
