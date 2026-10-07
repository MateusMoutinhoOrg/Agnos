package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeFlagAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_flag"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeFlagAction.RemoveFlag(sandbox, props.Path, entries.Target, entries.Name)
	if remove_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
