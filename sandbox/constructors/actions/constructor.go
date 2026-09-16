package actions

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	actions "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions"
)

// Constructor fills Sandbox.Actions, building it with the
// NewActions of sandbox/internal/actions. sandbox/new.go calls it
// once, along with the Constructor of every other package under
// sandbox/constructors/.
//
// Written once by `agnos build` and then yours: wrap the
// implementation, decorate the contract, or build a different one entirely.
// No build rewrites this file once it is there.
func Constructor(sandbox *api.Sandbox) {
	sandbox.Actions = actions.NewActions(sandbox)
}
