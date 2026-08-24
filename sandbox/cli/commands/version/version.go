package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

func NewCommand(sandbox *api.SandBox) api.CliCommand {
	return api.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		ArgsList:              []api.CliArg{},
		FlagsList:             []api.Cliflag{},
		Description:           "Returns version of program",
		Examples: []string{
			sandbox.ProjectName + " --version",
			sandbox.ProjectName + " version",
		},
		Handler: CommandHandler,
	}
}

func CommandHandler(d deps.Deps, argsRetriver *api.ArgsRetriver, flagsRetriver *api.FlagsRetriver) int {
	d.Printf("%s\n", "version")
	return api.ExitOk
}
