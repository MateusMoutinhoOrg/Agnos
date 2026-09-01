package smartio

import (
	"strings"

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

// isPendingRemoval checks if a path (or any of its parents) has been
// scheduled for removal in the current transaction.
func isPendingRemoval(io *SmartIO, path string) bool {
	for _, removed := range io.PendingRemoveDirs {
		if path == removed || strings.HasPrefix(path, removed+"/") {
			return true
		}
	}
	return false
}

// isPendingCreate checks if a path has been scheduled for creation
// in the current transaction.
func isPendingCreate(io *SmartIO, path string) bool {
	for _, created := range io.PendingCreateDirs {
		if path == created {
			return true
		}
	}
	return false
}

// filterPendingRemoved removes entries that are under a pending removal directory.
func filterPendingRemoved(io *SmartIO, paths []string) []string {
	var result []string
	for _, p := range paths {
		if !isPendingRemoval(io, p) {
			result = append(result, p)
		}
	}
	return result
}
