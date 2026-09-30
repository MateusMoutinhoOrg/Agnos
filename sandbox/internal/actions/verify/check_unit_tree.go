package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// checkUnitTree enforces what makes the tree under root readable by marker
// alone, the way CheckCommands and CheckRoutes read it: a unit's name is
// unique across every folder — the generated dispatch imports each under it —
// and a folder holding one of files without a marker beside them is a unit
// that lost its declaration, which no build collects and no verb addresses.
func checkUnitTree(sandbox *api.Sandbox, io *smartio.SmartIO, root string, marker string, kind string, files []string) []string {
	var violations []string

	units := utils.FindUnitDirs(sandbox, io, root, marker)
	seen := map[string]string{}
	for _, unit := range units {
		if other, taken := seen[unit.Name]; taken {
			violations = append(violations, unit.Dir+" declares the "+kind+" "+unit.Name+", which "+other+
				" already declares: a "+kind+" name is unique across every folder")
			continue
		}
		seen[unit.Name] = unit.Dir
	}

	for _, dir := range io.ListDirsRecursively(root) {
		if io.IsFile(dir + "/" + marker) {
			continue
		}
		for _, file := range files {
			if io.IsFile(dir + "/" + file) {
				violations = append(violations, dir+" holds "+file+" but no "+marker+
					": a directory is a "+kind+" by its "+marker+", so no build collects this one")
				break
			}
		}
	}

	return violations
}
