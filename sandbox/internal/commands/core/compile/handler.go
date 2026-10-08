package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	compile_error := compileAction.Compile(sandbox, api.CompileProps{
		Path:    props.Path,
		Targets: input.Target,
	})

	if compile_error != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", compile_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
