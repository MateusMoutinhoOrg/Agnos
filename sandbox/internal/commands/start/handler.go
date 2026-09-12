package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	var module *string
	if entries.Module != "" {
		modVal := entries.Module
		module = &modVal
	}

	if !sandbox.Deps.Iodeps.Exist(entries.Path+"/go.mod") && module == nil {
		{
			sandbox.Deps.Std.Error("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return api.ExitUsage
	}

	start_error := startAction.Start(sandbox, api.StartProps{
		Path:        entries.Path,
		ProjectName: entries.ProjectName,
		Module:      module,
		Force:       entries.Force,
	})

	if start_error != nil {
		sandbox.Deps.Std.Error("%s\n", start_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
