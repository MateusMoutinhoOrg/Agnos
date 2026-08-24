package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func CliMainFactory(sandbox *api.SandBox) func(args []string) int {

	return func(args []string) int {
		if len(args) == 0 {
			sandbox.Deps.Printf("%s", config.Usages)
			return api.ExitUsage
		}

		verb := sandbox.Deps.VerbLib
		action, error := verb.GetNextStringArg()

		for _, command := range sandbox.Commands {
			if contains(command.ValidStartIdentifiers, action) {

				return command.Handler(sandbox.Deps)
			}
		}

		sandbox.Deps.Printf("Unknown Command!")
		return api.ExitUsage

	}
}
