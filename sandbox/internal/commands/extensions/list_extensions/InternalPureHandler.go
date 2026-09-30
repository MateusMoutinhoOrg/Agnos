package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listExtensionsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_extensions"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	extensions, list_error := listExtensionsAction.ListExtensions(sandbox, props.Path)

	if list_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", list_error.Error())
	}

	for _, extension := range extensions {
		response.Printf("%-18s %-4s %s\n", extension.Name, stateMark(extension), extension.Help)
	}
	response.SetStatus(api.ExitOk)
	return nil
}

// stateMark is the middle column: whether agnos generates this mechanic for
// the project, or leaves it to be written by hand.
func stateMark(extension api.ExtensionInfo) string {
	if extension.Enabled {
		return "on"
	}
	return "off"
}
