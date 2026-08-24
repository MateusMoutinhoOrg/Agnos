package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
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

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {
	sandbox.Deps.Printf("Version:%s\n", sandbox.Version)
	return api.ExitOk
}
