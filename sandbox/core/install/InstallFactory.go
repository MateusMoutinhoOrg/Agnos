package install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// InstallFactory returns the closure that fills core.CoreApi.Install, adding
// the requested extension to the project in props.Path.
func InstallFactory(sandbox *sandbox.SandBox) func(props core.InstallProps) error {
	return func(props core.InstallProps) error {
		return nil
	}
}
