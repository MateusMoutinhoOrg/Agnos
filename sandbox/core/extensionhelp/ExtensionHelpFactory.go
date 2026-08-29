package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
)

// ExtensionHelpFactory returns the closure that fills
// lib.CoreApi.ExtensionHelp, printing the help of one extension.
func ExtensionHelpFactory(sandbox *lib.SandBox) func(props lib.ExtensionHelpProps) error {
	return func(props lib.ExtensionHelpProps) error {
		return nil
	}
}
