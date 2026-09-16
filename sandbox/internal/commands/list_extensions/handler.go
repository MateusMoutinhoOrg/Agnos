package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listExtensionsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_extensions"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	extensions, list_error := listExtensionsAction.ListExtensions(sandbox, command.GetString("path"))

	if list_error != nil {
		sandbox.Deps.Std.Error("%s\n", list_error.Error())
		return api.ExitFailure
	}

	for _, extension := range extensions {
		sandbox.Deps.Std.Printf("%-18s %-4s %s\n", extension.Name, stateMark(extension), extension.Help)
	}
	return api.ExitOk
}

// stateMark is the middle column: whether agnos generates this mechanic for
// the project, or leaves it to be written by hand.
func stateMark(extension api.ExtensionInfo) string {
	if extension.Enabled {
		return "on"
	}
	return "off"
}
