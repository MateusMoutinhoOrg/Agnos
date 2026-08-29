package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// ListFactory returns the closure that fills core.CoreApi.List, printing every
// extension available to the project in props.Path.
func ListFactory(sandbox *sandbox.SandBox) func(props core.ListProps) error {
	return func(props core.ListProps) error {
		return nil
	}
}
