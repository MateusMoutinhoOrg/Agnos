package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	set_error := setDepAction.SetDep(sandbox, api.SetDepProps{
		Path:            props.Path,
		Dep:             entries.Dep,
		Version:         entries.Version,
		RemoteAvailable: entries.RemoteAvailable,
	})

	if set_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
