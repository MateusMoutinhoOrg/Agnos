package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	renameRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/rename_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	rename_error := renameRouteAction.RenameRoute(sandbox, api.RenameRouteProps{
		Path:  props.Path,
		Route: input.Route,
		Name:  input.Name,
		// --dir declares no default, so "stay where it is" is no value at
		// all, and the top is spelled "/".
		Dir:    input.Dir,
		HasDir: input.Dir != "",
	})
	if rename_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", rename_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
