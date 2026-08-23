package sandbox

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

func SandBoxMainFactory(l *api.SandBox) func(args []string) int {
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

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.SandBox {
	l := api.SandBox{Deps: d}
	l.Sandboxmain = SandBoxMainFactory(&l)
	return l
}
