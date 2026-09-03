package set_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetCommand rewrites the command-level keys of
// sandbox/internal/commands/<command>/entries.yaml (help, category,
// long-description, hidden, identifiers, examples), then runs build.
func SetCommand(deps *deps.Deps, props api.CommandProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := SetCommandInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
