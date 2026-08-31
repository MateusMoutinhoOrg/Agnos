package smartio

import "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"

func ListFilesRecursively(deps *deps.Deps, io *SmartIO, path string) []string {
	p, err := processInputPath(deps, io, path)
	if err != nil {
		return nil
	}
	return filterIgnored(io, deps.IoLib.ListFilesRecursively(p))
}
