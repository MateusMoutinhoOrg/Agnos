package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	startAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/start"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	var module *string
	if input.Module != "" {
		modVal := input.Module
		module = &modVal
	}

	if !sandbox.Deps.IoDeps.Exists(props.Path+"/go.mod") && module == nil {
		{
			response.Eprintf("the module flag (--module) is required when there is no go.mod in the path\n")
		}
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitUsage, "", "")
	}

	start_error := startAction.Start(sandbox, api.StartProps{
		Path:        props.Path,
		ProjectName: input.ProjectName,
		Module:      module,
		Force:       input.Force,
	})

	if start_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", start_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
