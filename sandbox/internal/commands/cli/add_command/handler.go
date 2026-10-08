package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addCommandAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_command"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addCommandAction.AddCommand(sandbox, api.AddCommandProps{
		Path:              props.Path,
		Name:              input.Name,
		Trigger:           input.Trigger,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Pattern:           input.Pattern,
		Middleware:        input.Middleware,
		Priority:          input.Priority,
		HasPriority:       input.Priority >= 0,
		Before:            input.Before,
		After:             input.After,
		Summary:           input.Summary,
		Category:          input.Category,
		Dir:               input.Dir,
	})
	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
