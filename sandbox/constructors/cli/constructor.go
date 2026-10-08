package cli

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	cli "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/cli"
)

// Constructor fills Sandbox.Cli, building it with the
// NewCli of sandbox/internal/generated/cli. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `agnos build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.Cli = cli.NewCli(sandbox)
}
