package set_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	// --priority and --segments declare no default, so a value typed as 0
	// is told apart from none by whether anything was bound at all.
	set_error := setRouteAction.SetRoute(sandbox, api.RouteProps{
		Priority:        entries.Priority,
		HasPriority:     entries.Priority >= 0,
		Before:          entries.Before,
		After:           entries.After,
		Segments:        entries.Segments,
		HasSegments:     entries.Segments >= 0,
		Clear:           entries.Clear,
		Path:            props.Path,
		Route:           entries.Route,
		Methods:         entries.Method,
		ResponseType:    entries.ResponseType,
		Help:            entries.Help,
		Category:        entries.Category,
		LongDescription: entries.LongDescription,
		Hidden:          entries.Hidden,
		Visible:         entries.Visible,
		Examples:        entries.Example,
	})
	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
