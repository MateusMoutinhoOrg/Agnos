package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)
func NewCommand(sandbox *api.SandBox) api.CliCommand{
	return api.CliCommand{
		ValidStartIdentifiers: []string{"version", "--version"},
		Flags:                 []api.Cliflag{},              // TODO: flags
		Description:           "Returns version of program", // TODO: description
		Examples:              []string{
			sandbox.Deps.ProjectName + " --version",
			sandbox.Deps.ProjectName + " version"
		}, 
		Handler:               Command,
	}
}

func Command(sandbox *api.SandBox, fr api.FlagsRetriver) int {
	sandbox.Deps.Printf("%s\n")
	return api.ExitOk
}
