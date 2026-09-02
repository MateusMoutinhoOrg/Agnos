package version

import (
	api "{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
	"{{.Module}}/sandbox/internal/config"
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
	deps.Std.Printf("Version:%s\n", config.Version)
	return api.ExitOk
}
