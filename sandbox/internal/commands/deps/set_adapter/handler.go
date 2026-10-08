package set_adapter

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	setAdapterAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/set_adapter"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	set_error := setAdapterAction.SetAdapter(sandbox, api.SetAdapterProps{
		Path:    props.Path,
		Dep:     input.Dep,
		Adapter: input.Adapter,
		Binding: input.Binding,
	})

	if set_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", set_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
