package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// adapterlistDir is the tree `add-dep` renders the implementation half
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
func CheckAdapterlist(sandbox *api.Sandbox, io *smartio.SmartIO, module string) []string {
	var violations []string

	if !io.IsDir(adapterlistDir) {
		return violations
	}

	for _, dir := range io.ListDirs(adapterlistDir) {
		adapter := lastSegment(sandbox, dir)

		for _, asset := range io.ListFilesRecursively(dir) {
			relative := sandbox.Deps.Stringsdeps.TrimPrefix(asset, dir+"/")

			target := relative
			if relative == utils.AdapterConfFile {
				target = utils.AdapterConfPath(adapter)
			}

			violations = append(violations, checkCatalogAsset(sandbox, io, asset, target, module)...)
		}
	}

	return violations
}
