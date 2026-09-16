package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	disableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/disable_extension"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	disable_error := disableExtensionAction.DisableExtension(sandbox, command.GetString("path"), command.GetString("name"))

	if disable_error != nil {
		sandbox.Deps.Std.Error("%s\n", disable_error.Error())
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("%s is off; what it wrote is yours now\n", command.GetString("name"))
	return api.ExitOk
}
