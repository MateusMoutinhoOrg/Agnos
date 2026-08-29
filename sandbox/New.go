package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/lib"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core"
)

// New builds the lib.SandBox entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to lib.SandBox means adding its factory call here.
func New(d deps.Deps) lib.SandBox {
	sandbox := lib.SandBox{Deps: d}
	// Config runs first: the command definitions interpolate ProjectName
	// into their examples when they are built.
	Config(&sandbox)
	cli.SetCliMethods(&sandbox)
	core.SetCoreMethods(&sandbox)
	return sandbox
}
