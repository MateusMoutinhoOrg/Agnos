package binds

// This file is an illustrative sample, not part of the build.

import (
	api "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	deps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	depListAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/dep_list"
	verifyAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/verify"
)

// ActionsBind fills every field of sandbox.Actions, forwarding to the action
// packages with deps first. It holds no logic of its own.
func ActionsBind(deps *deps.Deps, sandbox *api.Sandbox) {
	sandbox.Actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(deps, props)
	}
	sandbox.Actions.Verify = func(path string) error {
		return verifyAction.Verify(deps, path)
	}
	sandbox.Actions.DepList = func(path string) ([]string, error) {
		return depListAction.DepList(deps, path)
	}
}
