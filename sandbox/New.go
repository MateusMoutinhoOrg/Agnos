package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/cli"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/core"
)

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.SandBox {
	sandbox := api.SandBox{Deps: d}
	cli.SetCliMethods(&sandbox)
	core.SetCoreMethods(&sandbox)
	Config(&sandbox)
	return sandbox
}
