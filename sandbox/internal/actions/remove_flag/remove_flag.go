package remove_flag

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveFlag drops one flag declaration from
// sandbox/internal/commands/<category>/<command>/command.yaml, then runs build so the
// generated.new.go forgets it.
func RemoveFlag(sandbox *api.Sandbox, props api.RemoveFlagProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveFlagInternal(sandbox, io, props.Command, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
