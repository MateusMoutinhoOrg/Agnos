package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

const (
	ExitOk      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

func CliMain(sandbox *api.Sandbox, args []string) int {

	verb := sandbox.Deps.Argvdeps.New(args)

	action, err := verb.GetNextStringArg()
	if err != nil {
		return ExitUsage
	}

	if action == "--version" || action == "-v" {
		sandbox.Deps.Std.Printf("%s: %s\n", sandbox.ProjectName, sandbox.Version)
		return ExitOk
	}

	return 0
}
