package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/commandprops"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cliio"
)

func InternalPureHandler(sandbox *api.Sandbox, props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error {
	compile_error := compileAction.Compile(sandbox, api.CompileProps{
		Path:    props.Path,
		Targets: entries.Target,
	})

	if compile_error != nil {
		return cliio.Fail(sandbox, api.ExitFailure, "", compile_error.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
