package set_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// SetArg rewrites one arg of sandbox/internal/commands/<command>/command.yaml,
// then runs build so the command's generated.new.go and generated.input.go follow it.
func SetArg(sandbox *api.Sandbox, props api.SetArgProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetArgInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
