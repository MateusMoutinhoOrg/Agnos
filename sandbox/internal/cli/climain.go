package cli

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func CliMain(deps *deps.Deps, args []string) int {

	verblib := deps.NewVerbLib(args)

	value, _ := verblib.GetNextStringArg()
	println(value)

	return 0
}
