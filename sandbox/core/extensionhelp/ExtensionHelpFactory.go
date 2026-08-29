package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/core"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib/sandbox"
)

// ExtensionHelpFactory returns the closure that fills
// core.CoreApi.ExtensionHelp, printing the help of one extension.
func ExtensionHelpFactory(sandbox *sandbox.SandBox) func(props core.ExtensionHelpProps) error {
	return func(props core.ExtensionHelpProps) error {
		return nil
	}
}
