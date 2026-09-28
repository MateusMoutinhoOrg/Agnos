package rename_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RenameCommand moves one command package to a new name, then runs build so
// the generated files and the dispatch follow it there.
func RenameCommand(sandbox *api.Sandbox, props api.RenameCommandProps) error {
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RenameCommandInternal(sandbox, io, props); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
