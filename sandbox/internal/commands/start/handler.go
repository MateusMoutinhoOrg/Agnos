package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	var module *string
	if command.GetString("module") != "" {
		modVal := command.GetString("module")
		module = &modVal
	}

	if !sandbox.Deps.Iodeps.Exist(command.GetString("path")+"/go.mod") && module == nil {
		{
			sandbox.Deps.Std.Error("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return api.ExitUsage
	}

	start_error := startAction.Start(sandbox, api.StartProps{
		Path:        command.GetString("path"),
		ProjectName: command.GetString("project-name"),
		Module:      module,
		Force:       command.GetBool("force"),
	})

	if start_error != nil {
		sandbox.Deps.Std.Error("%s\n", start_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
