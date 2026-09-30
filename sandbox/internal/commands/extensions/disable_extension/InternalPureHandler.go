package disable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	disableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/disable_extension"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	disable_error := disableExtensionAction.DisableExtension(sandbox, props.Path, entries.Name)

	if disable_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", disable_error.Error())
	}

	response.Printf("%s is off; what it wrote is yours now\n", entries.Name)
	response.SetStatus(api.ExitOk)
	return nil
}
