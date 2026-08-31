package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

	actions "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	actions.BindMethods(&self, *deps)
	return &self
}
