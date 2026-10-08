package set_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setDepAction.SetDep(sandbox, api.SetDepProps{
		Path:          props.Path,
		Dep:           input.Name,
		Version:       input.Version,
		RemoteBinding: input.RemoteBinding,
	})

	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
