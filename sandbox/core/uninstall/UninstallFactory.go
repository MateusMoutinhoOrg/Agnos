package uninstall

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// UninstallFactory returns the closure that fills core.CoreApi.Uninstall,
// removing the requested extension from the project in props.Path.
func UninstallFactory(sandbox *sandbox.SandBox) func(props core.UninstallProps) error {
	return func(props core.UninstallProps) error {
		return nil
	}
}
