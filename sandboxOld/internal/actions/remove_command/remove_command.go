package remove_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemoveCommand deletes the whole sandbox/internal/commands/<name>/ package,
// then runs build so climain.go and help stop dispatching to it.
func RemoveCommand(sandbox *api.Sandbox, path string, name string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := RemoveCommandInternal(sandbox, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
