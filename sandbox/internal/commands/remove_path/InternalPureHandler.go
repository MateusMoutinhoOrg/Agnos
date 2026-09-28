package remove_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePathAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_path"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	err := removePathAction.RemovePath(sandbox, props.Path, entries.Route, entries.Id)
	if err != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
