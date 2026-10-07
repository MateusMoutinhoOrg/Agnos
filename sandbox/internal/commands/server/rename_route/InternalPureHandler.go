package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	renameRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	rename_error := renameRouteAction.RenameRoute(sandbox, api.RenameRouteProps{
		Path:  props.Path,
		Route: entries.Route,
		Name:  entries.Name,
		// --dir declares no default, so "stay where it is" is no value at
		// all, and the top is spelled "/".
		Dir:    entries.Dir,
		HasDir: entries.Dir != "",
	})
	if rename_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", rename_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
