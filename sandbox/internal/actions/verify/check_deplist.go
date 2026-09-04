package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// deplistDir is the tree `dep-install` renders from: one directory per
// installable dep, each mirroring the layout of the project it is rendered
// into.
const deplistDir = "assets/deplist"

// moduleVar is the only template variable a deplist asset may use, so
// rendering one is a plain substitution of the target module path.
const moduleVar = "{{.Module}}"

// CheckDeplist enforces that every installable dep stays byte-identical to the
// copy this project runs on: assets/deplist/<dep>/<path> rendered with this
// module must equal <path> whenever that file exists here. Nothing else keeps
// the two in step, so a contract that gains a field in the project would
// otherwise go on being installed without it.
//
// A project with no assets/deplist has nothing to check.
func CheckDeplist(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	if !io.IsDir(deplistDir) {
		return violations
	}

	for _, dep := range io.ListDirs(deplistDir) {
		for _, asset := range io.ListFilesRecursively(dep) {
			target := strings.TrimPrefix(asset, dep+"/")
			if !io.IsFile(target) {
				continue
			}
			violations = append(violations, checkDeplistAsset(io, asset, target, module)...)
		}
	}

	return violations
}

// checkDeplistAsset compares one rendered asset with the file it installs
// over. A target that is absent from this project is not a violation: a dep
// this project does not use has nothing here to drift from.
func checkDeplistAsset(io *smartio.SmartIO, asset string, target string, module string) []string {
	source, err := io.ReadFile(asset)
	if err != nil {
		return []string{asset + " could not be read"}
	}

	installed, err := io.ReadFile(target)
	if err != nil {
		return []string{target + " could not be read"}
	}

	if strings.ReplaceAll(string(source), moduleVar, module) == string(installed) {
		return nil
	}

	return []string{asset + " has drifted from " + target +
		" (an installable dep must render byte-for-byte to the copy this project runs on)"}
}
