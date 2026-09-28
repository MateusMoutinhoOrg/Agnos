package list_commands

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListCommandsInternal is every declared command as one line, in run order:
// the rung, whether it is strict or a middleware, the pattern and the name —
// the chain a command line walks until one command answers.
func ListCommandsInternal(sandbox *api.Sandbox, io *smartio.SmartIO) ([]string, error) {
	if err := utils.RequireProject(sandbox, io); err != nil {
		return nil, err
	}
	if !io.IsDir(utils.CommandsDir) {
		return nil, sandbox.Deps.Std.Errorf("the project has no cli layer: run cli-init first")
	}

	chain, err := utils.LoadCommandChain(sandbox, io)
	if err != nil {
		return nil, err
	}

	lines := []string{sandbox.Deps.Std.Sprintf("%-4s %-9s %-10s %-36s %s", "#", "priority", "kind", "pattern", "command")}
	for index, entry := range chain {
		kind := "command"
		if !entry.Conf.Strict {
			kind = "middleware"
		}
		lines = append(lines, sandbox.Deps.Std.Sprintf("%-4d %-9d %-10s %-36s %s",
			index+1, entry.Conf.Priority, kind, entry.Conf.Pattern(), utils.CommandIdentifier(sandbox, entry.Name)))
	}
	return lines, nil
}
