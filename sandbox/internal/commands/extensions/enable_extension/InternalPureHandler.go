package enable_extension

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	enableExtensionAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/enable_extension"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	enable_error := enableExtensionAction.EnableExtension(sandbox, props.Path, entries.Name)

	if enable_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", enable_error.Error())
	}

	response.Printf("%s is on\n", entries.Name)
	response.SetStatus(api.ExitOk)
	return nil
}
