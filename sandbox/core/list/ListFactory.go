package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// ListFactory returns the closure that fills api.CoreApi.List, printing every
// extension available to the project in props.Path.
func ListFactory(sandbox *api.SandBox) func(props api.ListProps) error {
	return func(props api.ListProps) error {
		return nil
	}
}
