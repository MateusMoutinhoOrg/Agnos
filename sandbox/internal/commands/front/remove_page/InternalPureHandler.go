package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removePageAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_page"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removePageAction.RemovePage(sandbox, props.Path, entries.Name)

	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
