package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

	actions "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions"
	cli "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/cli"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	actions.Export(deps, &self)
	cli.Export(deps, &self)
	return &self
}
