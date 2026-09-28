package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetCommand rewrites the command-level keys of
// sandbox/internal/commands/<command>/command.yaml (help, category,
// long-description, priority, segments, strict, hidden, identifiers,
// examples), then runs build.
func SetCommand(sandbox *api.Sandbox, props api.SetCommandProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := SetCommandInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
