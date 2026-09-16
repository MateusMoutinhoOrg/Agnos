package enable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	enableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/enable_extension"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	enable_error := enableExtensionAction.EnableExtension(sandbox, command.GetString("path"), command.GetString("name"))

	if enable_error != nil {
		sandbox.Deps.Std.Error("%s\n", enable_error.Error())
		return api.ExitFailure
	}

	sandbox.Deps.Std.Printf("%s is on\n", command.GetString("name"))
	return api.ExitOk
}
