package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

	//Exposed Methods
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions"
)

func New(deps *deps.Deps) *api.Sandbox {
	api := api.Sandbox{}
	api.Build = func(props api.BuildProps) error {
		return actions.Build(deps, props)
	}
	return &api
}
