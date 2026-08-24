package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		Description:           "Returns version of program",
		Examples: []string{
			sandbox.ProjectName + " --version",
			sandbox.ProjectName + " version",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(d deps.Deps, entries api.CliEntrys) int {
	d.Printf("%s\n", "version")
	return api.ExitOk
}
