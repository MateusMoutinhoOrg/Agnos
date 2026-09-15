package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
)

func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	compile_error := compileAction.Compile(sandbox, api.CompileProps{
		Path:    command.GetString("path"),
		Targets: command.GetStrings("target"),
	})

	if compile_error != nil {
		sandbox.Deps.Std.Error("%s\n", compile_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
