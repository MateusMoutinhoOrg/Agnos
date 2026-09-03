package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	compileAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/compile"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	compile_error := compileAction.Compile(deps, api.CompileProps{
		Path:    entries.Path,
		Targets: entries.Target,
	})

	if compile_error != nil {
		deps.Std.Error("%s\n", compile_error.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
