package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	actionscostructor "github.com/MateusMoutinhoOrg/Agnos/sandbox/constructors/actions"
	clicostructor "github.com/MateusMoutinhoOrg/Agnos/sandbox/constructors/cli"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	clicostructor.Constructor(&self)
	actionscostructor.Constructor(&self)

	return &self
}
