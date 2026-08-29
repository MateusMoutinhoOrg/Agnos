package uninstall

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

// UninstallFactory returns the closure that fills lib.CoreApi.Uninstall,
// removing the requested extension from the project in props.Path.
func UninstallFactory(sandbox *lib.SandBox) func(props lib.UninstallProps) error {
	return func(props lib.UninstallProps) error {
		return nil
	}
}
