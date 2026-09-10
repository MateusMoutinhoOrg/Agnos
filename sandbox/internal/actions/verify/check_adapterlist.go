package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// adapterlistDir is the tree `dep-install` renders the implementation half
// from: one directory per installable adapter, each holding its adapter.yaml
// beside a mirror of the layout it is rendered into.
const adapterlistDir = "assets/" + utils.AdapterlistGroup

// CheckAdapterlist is CheckDeplist for the other half of the catalog: every
// installable adapter must stay byte-identical to the copy this project runs
// on. The adapter.yaml is the one asset that does not sit at the path it
// installs to — it declares the package rather than belonging to the mirror —
// so it is compared against the copy install writes into the package.
//
// A project with no assets/adapterlist has nothing to check.
func CheckAdapterlist(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	if !io.IsDir(adapterlistDir) {
		return violations
	}

	for _, dir := range io.ListDirs(adapterlistDir) {
		adapter := lastSegment(deps, dir)

		for _, asset := range io.ListFilesRecursively(dir) {
			relative := deps.Stringsdeps.TrimPrefix(asset, dir+"/")

			target := relative
			if relative == utils.AdapterConfFile {
				target = utils.AdapterConfPath(adapter)
			}

			violations = append(violations, checkCatalogAsset(deps, io, asset, target, module)...)
		}
	}

	return violations
}
