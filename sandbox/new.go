package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	actions "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/cli"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/extensions"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	self.Actions = actions.NewActions(&self)
	self.Cli = cli.NewCli(&self)
	self.Extensions = extensions.NewExtensions(&self)
	return &self
}
