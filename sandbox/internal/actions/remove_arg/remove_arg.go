package remove_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveArg drops one positional arg declaration from
// sandbox/internal/commands/<category>/<command>/command.yaml, then runs build so the
// generated.new.go forgets it.
func RemoveArg(sandbox *api.Sandbox, props api.RemoveArgProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveArgInternal(sandbox, io, props.Command, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
