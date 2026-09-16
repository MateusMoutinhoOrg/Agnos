package actions

import (
	api "github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	actions "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions"
)

func Constructor(sandbox *api.Sandbox) {

	sandbox.Actions = actions.NewActions(sandbox)
}
