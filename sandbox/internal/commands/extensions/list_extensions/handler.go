package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listExtensionsAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_extensions"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	extensions, list_error := listExtensionsAction.ListExtensions(sandbox, api.ListExtensionsProps{Path: props.Path})

	if list_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", list_error.Error())
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
