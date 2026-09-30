package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	var module *string
	if entries.Module != "" {
		modVal := entries.Module
		module = &modVal
	}

	if !sandbox.Deps.Iodeps.Exist(props.Path+"/go.mod") && module == nil {
		{
			response.Error("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return cliio.Fail(sandbox, api.ExitUsage, "", "")
	}

	start_error := startAction.Start(sandbox, api.StartProps{
		Path:        props.Path,
		ProjectName: entries.ProjectName,
		Module:      module,
		Force:       entries.Force,
	})

	if start_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", start_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
