package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// deplistDir is the tree `add-dep` renders the contract half from: one
// directory per installable dep, each holding its dep.yaml beside a mirror of
// the layout it is rendered into.
const deplistDir = "assets/" + utils.DeplistGroup

// moduleVar is the only template variable a catalog asset may use, so
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
			relative := deps.Stringsdeps.TrimPrefix(asset, dep+"/")
			if relative == utils.DepConfFile {
				continue
			}
			violations = append(violations, checkCatalogAsset(deps, io, asset, relative, module)...)
		}
	}

	return violations
}

// checkCatalogAsset compares one rendered catalog asset with the file it
// installs over. A target that is absent from this project is not a violation:
// a dep or an adapter this project does not use has nothing here to drift from.
func checkCatalogAsset(deps *deps.Deps, io *smartio.SmartIO, asset string, target string, module string) []string {
	if !io.IsFile(target) {
		return nil
	}

	source, err := io.ReadFile(asset)
	if err != nil {
		return []string{asset + " could not be read"}
	}

	installed, err := io.ReadFile(target)
	if err != nil {
		return []string{target + " could not be read"}
	}

	if deps.Stringsdeps.ReplaceAll(string(source), moduleVar, module) == string(installed) {
		return nil
	}

	return []string{asset + " has drifted from " + target +
		" (an installable dep must render byte-for-byte to the copy this project runs on)"}
}
