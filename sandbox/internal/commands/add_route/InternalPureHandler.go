package add_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	add_error := addRouteAction.AddRoute(sandbox, api.AddRouteProps{
		Path:              props.Path,
		Name:              entries.Name,
		Methods:           entries.Method,
		Trigger:           entries.Trigger,
		TriggerType:       entries.TriggerType,
		TriggerNegate:     entries.TriggerNegate,
		TriggerIgnoreCase: entries.TriggerIgnoreCase,
		Pattern:           entries.Pattern,
		Middleware:        entries.Middleware,
		Priority:          entries.Priority,
		// --priority declares no default, so a rung typed as 0 is told
		// apart from none by whether anything was bound at all.
		HasPriority:  entries.Priority >= 0,
		Before:       entries.Before,
		After:        entries.After,
		ResponseType: entries.ResponseType,
		Help:         entries.Help,
		Category:     entries.Category,
	})

	if add_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", add_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
