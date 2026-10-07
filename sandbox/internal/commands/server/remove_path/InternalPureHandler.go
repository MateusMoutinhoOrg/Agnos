package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := removePathAction.RemovePath(sandbox, props.Path, entries.Route, entries.Id)
	if err != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
