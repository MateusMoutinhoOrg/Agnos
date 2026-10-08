package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	// --priority and --segments declare no default, so a value typed as 0
	// is told apart from none by whether anything was bound at all.
	set_error := setRouteAction.SetRoute(sandbox, api.SetRouteProps{
		Priority:     input.Priority,
		HasPriority:  input.Priority >= 0,
		Before:       input.Before,
		After:        input.After,
		Segments:     input.Segments,
		HasSegments:  input.Segments >= 0,
		Clear:        input.Clear,
		Path:         props.Path,
		Route:        input.Name,
		Methods:      input.Method,
		ResponseType: input.ResponseType,
		Summary:      input.Summary,
		Category:     input.Category,
		Description:  input.Description,
		Hidden:       input.Hidden,
		Visible:      input.Visible,
		Examples:     input.Example,
	})
	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
