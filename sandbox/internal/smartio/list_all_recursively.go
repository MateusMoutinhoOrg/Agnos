package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func ListAllRecursively(deps *deps.Deps, io *SmartIO, path string) []string {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil
	}
	return filterIgnored(io, deps.IoLib.ListAllRecursively(p))
}
