package list_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// ListCommands renders the chain in the order the dispatch runs it. Like the
// other readers it opens a StagedFS and never calls io.Persist, and it runs no
// follow-up build.
func ListCommands(sandbox *api.Sandbox, props api.ListCommandsProps) ([]string, error) {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	return ListCommandsInternal(sandbox, io)
}
