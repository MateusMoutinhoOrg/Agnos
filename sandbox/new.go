package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	actions "github.com/MateusMoutinhoOrg/Agnos/sandbox/constructors/actions"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/constructors/cli"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// New builds the whole library: one call per package under
// sandbox/constructors/, each filling the field of the Sandbox it owns. The
// list is the directories themselves, so a constructor written by hand is
// called exactly like a generated one — this file is rendered around what is
// there, never the other way round.
func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	actions.Constructor(&self)
	cli.Constructor(&self)

	return &self
}
