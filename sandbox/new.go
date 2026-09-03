package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"

	binds "github.com/MateusMoutinhoOrg/Agnos/sandbox/binds"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	binds.ActionsBind(deps, &self)
	binds.CliBind(deps, &self)

	return &self
}
