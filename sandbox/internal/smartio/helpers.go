package smartio

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

func processInputPath(deps *deps.Deps, io *SmartIO, path string) (string, error) {
	p := io.Replacers.Format(path)
	if io.Ignore.IsIgnorable(p) {
		return p, deps.Errorf("path %q is ignorable", p)
	}
	return p, nil
}

func filterIgnored(io *SmartIO, paths []string) []string {
	var result []string
	for _, p := range paths {
		if !io.Ignore.IsIgnorable(p) {
			result = append(result, p)
		}
	}
	return result
}
