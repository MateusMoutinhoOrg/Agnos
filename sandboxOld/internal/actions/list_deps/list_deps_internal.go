package list_deps

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListDepsInternal returns one row per dep of the embedded catalog, in listing
// order, saying which the project has installed and which adapters fill each
// one — the two halves the catalog was split into, answered together.
func ListDepsInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) ([]api.DepInfo, error) {
	sandbox.Deps.Std.Log("list-deps started with path %s \n", path)

	catalog, err := utils.CatalogDeps(sandbox)
	if err != nil {
		return nil, err
	}

	deplist := []api.DepInfo{}
	for _, name := range catalog {
		conf, err := utils.LoadCatalogDepConf(sandbox, name)
		if err != nil {
			return nil, err
		}

		deplist = append(deplist, api.DepInfo{
			Name:           name,
			Field:          conf.Field,
			Help:           conf.Help,
			DefaultAdapter: conf.DefaultAdapter,
			Installed:      io.IsDir(utils.ContractsDir + "/" + name),
			Adapters:       utils.AdaptersFillingDep(sandbox, io, name),
		})
	}

	return deplist, nil
}
