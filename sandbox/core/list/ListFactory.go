package list

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

// ListFactory returns the closure that fills lib.CoreApi.List, printing every
// extension available to the project in props.Path.
func ListFactory(sandbox *lib.SandBox) func(props lib.ListProps) error {
	return func(props lib.ListProps) error {
		return nil
	}
}
