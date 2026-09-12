package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetCommand rewrites the command-level keys of
// sandbox/internal/commands/<command>/entries.yaml (help, category,
// long-description, hidden, identifiers, examples), then runs build.
func SetCommand(sandbox *api.Sandbox, props api.CommandProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := SetCommandInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
