package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveArg drops one positional arg declaration from
// sandbox/internal/commands/<command>/entries.yaml, then runs build so the
// generated new.go forgets it.
func RemoveArg(sandbox *api.Sandbox, path string, command string, name string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := RemoveArgInternal(sandbox, io, command, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
