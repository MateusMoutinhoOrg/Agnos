package add_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	install_error := addAdapterAction.AddAdapter(sandbox, api.AddAdapterProps{
		Path:    props.Path,
		Adapter: input.Name,
		Binding: input.Binding,
	})

	if install_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", install_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
