package config

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	config "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/config"
)

// Constructor fills Sandbox.Config, building it with the
// NewConfig of sandbox/internal/generated/config. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `agnos build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.Config = config.NewConfig(sandbox)
}
