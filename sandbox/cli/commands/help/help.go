package help

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"help", "--help"},
		ArgsList: []api.CliArg{
			api.CliArg{
				Name:        "command",
				Description: "The command to get help for",
				Required:    false,
				Type:        api.CliTypeString,
				Size:        1,
			},
		},
		FlagsList:   []api.Cliflag{},
		Description: "Shows Help of a command",
		Examples: []string{
			sandbox.ProjectName + " --help",
			sandbox.ProjectName + " help",
			sandbox.ProjectName + " help <command>",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(d deps.Deps, argsRetriver *api.ArgsRetriver, flagsRetriver *api.FlagsRetriver) int {
	command := argsRetriver.GetStringArg("command", 0)

	d.Printf("%s\n", command)
	return api.ExitOk
}
