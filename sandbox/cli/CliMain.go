package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

func CliMainFactory(l *api.SandBox) func(args []string) int {

	return func(args []string) int {
		if len(args) == 0 {
			l.Deps.Printf("%s", config.Usages)
			return api.ExitUsage
		}

		verb := l.Deps.VerbLib
		if verb.IsPresent(config.HelpFlags) {
			l.Deps.Printf("%s\n", config.Usages)
			return api.ExitOk
		}
		if verb.IsPresent(config.VersionFlags) {
			l.Deps.Printf("%s\n", config.Version)
			return api.ExitOk
		}

		l.Deps.Printf("Unknown Command!")
		return api.ExitUsage

	}
}
