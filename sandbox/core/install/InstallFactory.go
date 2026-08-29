package install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

// InstallFactory returns the closure that fills lib.CoreApi.Install, adding
// the requested extension to the project in props.Path.
func InstallFactory(sandbox *lib.SandBox) func(props lib.InstallProps) error {
	return func(props lib.InstallProps) error {
		return nil
	}
}
