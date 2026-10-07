package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	install_error := addDepAction.AddDep(sandbox, api.AddDepProps{
		Path:            props.Path,
		Dep:             entries.Dep,
		Adapter:         entries.Adapter,
		As:              entries.As,
		RemoteAvailable: entries.RemoteAvailable,
	})

	if install_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", install_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
