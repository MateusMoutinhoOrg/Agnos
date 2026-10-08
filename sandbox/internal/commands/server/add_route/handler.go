package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	add_error := addRouteAction.AddRoute(sandbox, api.AddRouteProps{
		Path:              props.Path,
		Name:              input.Name,
		Methods:           input.Method,
		Trigger:           input.Trigger,
		TriggerType:       input.TriggerType,
		TriggerNegate:     input.TriggerNegate,
		TriggerIgnoreCase: input.TriggerIgnoreCase,
		Pattern:           input.Pattern,
		Middleware:        input.Middleware,
		Priority:          input.Priority,
		// --priority declares no default, so a rung typed as 0 is told
		// apart from none by whether anything was bound at all.
		HasPriority:  input.Priority >= 0,
		Before:       input.Before,
		After:        input.After,
		ResponseType: input.ResponseType,
		Summary:      input.Summary,
		Category:     input.Category,
		Dir:          input.Dir,
	})

	if add_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
