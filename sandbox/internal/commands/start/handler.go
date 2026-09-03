package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	startAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/start"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	var module *string
	if entries.Module != "" {
		modVal := entries.Module
		module = &modVal
	}

	if !deps.Iodeps.Exist(entries.Path+"/go.mod") && module == nil {
		{
			deps.Std.Error("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return api.ExitUsage
	}

	start_error := startAction.Start(deps, api.StartProps{
		Path:        entries.Path,
		ProjectName: entries.ProjectName,
		Module:      module,
		Force:       entries.Force,
	})

	if start_error != nil {
		deps.Std.Error("%s\n", start_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
