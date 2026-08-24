package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(deps deps.Deps) api.SandBox {
	sandbox := api.SandBox{Deps: depss}
	NewCli(&sandbox)
	return sandbox
}
