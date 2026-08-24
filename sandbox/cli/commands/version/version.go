package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

var Command = api.CliCommand{
	ValidStartIdentifiers: []string{"version", "--version"},
	Flags:                 []api.Cliflag{},       // TODO: flags
	Description:           "Returns version of ", // TODO: description
	Examples:              []string{"version"},   // TODO: examples
	Handler:               command_version,
}

func command_version(sandbox *api.SandBox, fr api.FlagsRetriver) int {
	sandbox.Deps.Printf("%s\n")
	return api.ExitOk
}
