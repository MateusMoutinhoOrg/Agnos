package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func Export(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(props api.BuildProps) error {
		return Build(deps, props)
	}
	sandbox.Actions.InstallExtension = func(props api.InstallProps) error {
		return InstallExtension(deps, props)
	}
	sandbox.Actions.RemoveExtension = func(props api.UninstallProps) error {
		return RemoveExtension(deps, props)
	}
	sandbox.Actions.ListAvaliableExtensions = func(props api.ListAvailableExtensionsProps) ([]string, error) {
		return ListAvailableExtensions(deps, props)
	}
}
