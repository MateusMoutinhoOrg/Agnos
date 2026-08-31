package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func New(deps *deps.Deps) *api.Api {
	api := api.Api{}

	return &api
}
