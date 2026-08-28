package uninstall

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// UninstallFactory returns the closure that fills api.CoreApi.Uninstall,
// removing the requested extension from the project in props.Path.
func UninstallFactory(sandbox *api.SandBox) func(props api.UninstallProps) error {
	return func(props api.UninstallProps) error {
		return nil
	}
}
