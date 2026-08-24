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
		if verb.IsPresent(config.HelpFlags) {
			sandbox.Deps.Printf("%s\n", config.Usages)
			return api.ExitOk
		}
		if verb.IsPresent(config.VersionFlags) {
			sandbox.Deps.Printf("%s\n", config.Version)
			return api.ExitOk
		}

		sandbox.Deps.Printf("Unknown Command!")
		return api.ExitUsage

	}
}
