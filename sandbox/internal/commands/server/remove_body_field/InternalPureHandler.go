package remove_body_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	removeBodyFieldAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/remove_body_field"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	remove_error := removeBodyFieldAction.RemoveBodyField(sandbox, props.Path, entries.Route, entries.Name)

	if remove_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", remove_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
