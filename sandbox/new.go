package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/cli"
	config "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	self.Cli = cli.NewCli(&self)
	self.Version = config.Version
	self.ProjectName = config.ProjectName
	return &self
}
