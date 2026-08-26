package help

import (
	"slices"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"help", "--help"},
		Args: []api.CliArg{
			api.CliArg{
				Id:              "command",
				Description:     "The command to get help for",
				RequiredType:    api.CliTypeString,
				RequiredMinSize: 0,
				RequiredMaxSize: 1,
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

func CommandHandler(sandbox *api.SandBox, entries api.CliEntrys) int {
	command := entries.GetArgById("command")
	if len(command.Values) == 0 {
		sandbox.Deps.Printf("average  help here ")
		return api.ExitOk
	}
	chosen := command.Values[0].String()
	for _, c := range sandbox.Commands {
		if slices.Contains(c.ValidStartIdentifiers, chosen) {
			sandbox.Deps.Printf("%s\n", c.Description)
			return api.ExitOk
		}
	}

	sandbox.Deps.Printf("Unknow Command: %s\n", chosen)
	return api.ExitOk
}
