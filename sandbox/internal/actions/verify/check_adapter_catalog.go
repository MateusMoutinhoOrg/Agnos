package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// adapterCatalogDir is the tree `add-dep` renders the implementation half
// from: one directory per installable adapter, each holding its adapter.yaml
// beside a mirror of the layout it is rendered into.
const adapterCatalogDir = "assets/" + utils.AdapterCatalogGroup

// CheckAdapterCatalog is CheckDepCatalog for the other half of the catalog: every
// installable adapter must stay byte-identical to the copy this project runs
// on. The adapter.yaml is the one asset that does not sit at the path it
// installs to — it declares the package rather than belonging to the mirror —
// so it is compared against the copy install writes into the package.
//
// A project with no assets/adapter-catalog has nothing to check.
func CheckAdapterCatalog(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string) []string {
	var violations []string

	if !io.IsDir(adapterCatalogDir) {
		return violations
	}

	for _, dir := range io.ListDirs(adapterCatalogDir) {
		adapter := lastSegment(sandbox, dir)

		for _, asset := range io.ListFilesRecursively(dir) {
			relative := sandbox.Deps.StringsDeps.TrimPrefix(asset, dir+"/")

			target := relative
			if relative == utils.AdapterConfFile {
				target = utils.AdapterConfPath(adapter)
			}

			violations = append(violations, checkCatalogAsset(sandbox, io, asset, target, module)...)
		}
	}

	return violations
}
