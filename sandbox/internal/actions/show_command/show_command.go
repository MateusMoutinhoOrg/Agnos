package show_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ShowCommand renders one command's whole declaration as the lines of a tree.
// Like the list actions it opens a StagedFS and never calls io.Persist, and it
// runs no follow-up build: reading a declaration changes nothing.
func ShowCommand(sandbox *api.Sandbox, props api.ShowCommandProps) ([]string, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ShowCommandInternal(sandbox, io, props.Name)
}
