package cli_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CliPurge removes from the project every file the "cli" asset group would
// have installed, then runs build as a follow-up step.
func CliPurge(deps *deps.Deps, path string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := CliPurgeInternal(deps, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, path)
}
