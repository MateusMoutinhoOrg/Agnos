package show_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// ShowCommand renders one command's whole declaration as the lines of a tree.
// Like the list actions it opens a SmartIO and never calls io.Persist, and it
// runs no follow-up build: reading a declaration changes nothing.
func ShowCommand(sandbox *api.Sandbox, path string, command string) ([]string, error) {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	return ShowCommandInternal(sandbox, io, command)
}
