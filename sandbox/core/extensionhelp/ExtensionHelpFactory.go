package extensionhelp

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// ExtensionHelpFactory returns the closure that fills
// api.CoreApi.ExtensionHelp, printing the help of one extension.
func ExtensionHelpFactory(sandbox *api.SandBox) func(props api.ExtensionHelpProps) error {
	return func(props api.ExtensionHelpProps) error {
		return nil
	}
}
