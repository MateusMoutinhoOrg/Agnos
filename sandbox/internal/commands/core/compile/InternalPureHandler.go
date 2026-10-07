package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	compile_error := compileAction.Compile(sandbox, api.CompileProps{
		Path:    props.Path,
		Targets: entries.Target,
	})

	if compile_error != nil {
		return sandbox.Deps.OpinatedAgnosCli.Fail(api.ExitFailure, "", compile_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
