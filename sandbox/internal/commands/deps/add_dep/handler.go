package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	install_error := addDepAction.AddDep(sandbox, api.AddDepProps{
		Path:          props.Path,
		Dep:           input.Name,
		Adapter:       input.Adapter,
		As:            input.As,
		RemoteBinding: input.RemoteBinding,
	})

	if install_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", install_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
