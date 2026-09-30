package show_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	showRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/show_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, show_error := showRouteAction.ShowRoute(sandbox, props.Path, entries.Route)

	if show_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", show_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
