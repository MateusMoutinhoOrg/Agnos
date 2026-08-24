package help

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"help", "--help"},
		Args: []api.CliArg{
			api.CliArg{
				Name:         "command",
				Description:  "The command to get help for",
				RequiredType: api.CliTypeString,
				RequiredSize: 1,
			},
		},

		Description: "Shows Help of a command",
		Examples: []string{
			sandbox.ProjectName + " --help",
			sandbox.ProjectName + " help",
			sandbox.ProjectName + " help <command>",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(d deps.Deps, entries api.CliEntrys) int {
	command := entries.GetArgByName("command")

	if len(command.Values) == 0 {
		d.Printf("No command provided\n")
		return api.ExitUsage
	}
	user_command := command.Values[0].String()

	return api.ExitOk
}
