package install

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// InstallFactory returns the closure that fills api.CoreApi.Install, adding
// the requested extension to the project in props.Path.
func InstallFactory(sandbox *api.SandBox) func(props api.InstallProps) error {
	return func(props api.InstallProps) error {
		return nil
	}
}
