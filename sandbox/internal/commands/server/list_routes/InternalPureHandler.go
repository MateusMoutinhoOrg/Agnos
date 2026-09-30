package list_routes

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	listRoutesAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/list_routes"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	lines, list_error := listRoutesAction.ListRoutes(sandbox, props.Path)
	if list_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", list_error.Error())
	}

	for _, line := range lines {
		response.Printf("%s\n", line)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
