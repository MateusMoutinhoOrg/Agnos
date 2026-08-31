package actions

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func ExportMethods(deps deps.Deps, sandbox *api.Sandbox) {
	sandbox.Build = func(props api.BuildProps) error {
		return Build(deps, props)
	}
	sandbox.InstallExtension = func(props api.InstallProps) error {
		return InstallExtension(deps, props)
	}
	sandbox.RemoveExtension = func(props api.UninstallProps) error {
		return RemoveExtension(deps, props)
	}
	sandbox.ListAvaliableExtensions = func(props api.ListAvailableExtensionsProps) ([]string, error) {
		return ListAvailableExtensions(deps, props)
	}
}
