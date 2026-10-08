package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// SetCommand rewrites the command-level keys of
// sandbox/internal/commands/<command>/command.yaml (summary, category,
// description, priority, segments, strict, hidden, identifiers, examples),
// then runs build.
func SetCommand(sandbox *api.Sandbox, props api.SetCommandProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetCommandInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
