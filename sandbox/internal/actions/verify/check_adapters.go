package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// adaptersAllowedDirs is the fixed set of sub-directories the adapters/ tree
// may contain.
var adaptersAllowedDirs = []string{"availables", "libs"}

// CheckAdapters enforces that adapters/ holds only the availables and libs
// directories.
func CheckAdapters(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir("adapters") {
		return violations
	}

	for _, dir := range io.ListDirs("adapters") {
		name := lastSegment(dir)
		if !contains(adaptersAllowedDirs, name) {
			violations = append(violations, "adapters/ contains unexpected directory "+name+
				" (allowed: availables, libs)")
		}
	}

	for _, file := range io.ListFiles("adapters") {
		violations = append(violations, "adapters/ contains unexpected file "+lastSegment(file)+
			" (adapters/ may hold only the availables and libs directories)")
	}

	return violations
}
