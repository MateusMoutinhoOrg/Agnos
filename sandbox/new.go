package sandbox

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

	//Exposed Methods
	actions "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions"
)

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{}
	self.Build = func(props api.BuildProps) error {
		return actions.Build(*deps, props)
	}
	self.InstallExtension = func(props api.InstallProps) error {
		return actions.InstallExtension(*deps, props)
	}
	self.RemoveExtension = func(props api.UninstallProps) error {
		return actions.RemoveExtension(*deps, props)
	}
	self.ListAvaliableExtensions = func(props api.ListAvailableExtensionsProps) ([]string, error) {
		return actions.ListAvailableExtensions(*deps, props)
	}
	return &self
}
