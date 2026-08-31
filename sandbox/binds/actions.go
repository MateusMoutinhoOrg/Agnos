package binds

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	actions "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions"
)

func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(path string) error {
		return actions.Build(deps, path)
	}
}
