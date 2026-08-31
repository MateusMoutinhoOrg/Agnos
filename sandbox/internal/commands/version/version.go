package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

func NewCommand(deps *deps.Deps, sandbox *api.Sandbox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		Category:              "Info",
		Description:           "Print the installed version",
		LongDescription:       "Prints the current version of the installed binary and exits.",
		Examples: []string{
			config.ProjectName + " --version",
			config.ProjectName + " version",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(deps *deps.Deps, entries api.CliEntrys) int {
	deps.Printf("Version:%s\n", config.Version)
	return api.ExitOk
}
