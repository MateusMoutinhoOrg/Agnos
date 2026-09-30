package deps_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	depsInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/deps_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	init_error := depsInitAction.DepsInit(sandbox, props.Path)

	if init_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", init_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
