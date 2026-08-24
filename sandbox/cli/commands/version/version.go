package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)
func NewCommand(sandbox *api.SandBox) api.CliCommand{
	return api.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		ArgsList: []api.CliArg{},
		FlagsList:                 []api.Cliflag{},             
		Description:           "Returns version of program", 
		Examples:              []string{
			sandbox.Deps.ProjectName + " --version",
			sandbox.Deps.ProjectName + " version"
		}, 
		Handler:               CommandHandler,
	}
}

func CommandHandler(sandbox *api.SandBox, argsRetriver api.ArgsRetriver, flagsRetriver api.FlagsRetriver) int {
	sandbox.Deps.Printf("%s\n")
	return api.ExitOk
}
