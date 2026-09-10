package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListAdaptersInternal returns one row per adapter the project can reach: every
// adapter of the embedded catalog, plus every one installed that the catalog
// does not have — a generated shim is only ever in the second group.
func ListAdaptersInternal(deps *deps.Deps, io *smartio.SmartIO, path string) ([]api.AdapterInfo, error) {
	deps.Std.Log("list-adapters started with path %s \n", path)

	catalog, err := utils.CatalogAdapters(deps)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	adapters := []api.AdapterInfo{}

	for _, name := range catalog {
		seen[name] = true
		conf, err := utils.LoadCatalogAdapterConf(deps, name)
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, adapterInfo(deps, io, name, conf.Dep, conf.Help, conf.Module, conf.Origin))
	}

	for _, name := range utils.InstalledAdapters(deps, io) {
		if seen[name] {
			continue
		}
		conf, err := utils.LoadAdapterConf(deps, io, name)
		if err != nil {
			continue
		}
		adapters = append(adapters, adapterInfo(deps, io, name, conf.Dep, conf.Help, conf.Module, conf.Origin))
	}

	return adapters, nil
}

// adapterInfo fills one row with what the project holds of that adapter: an
// installed package, and the availables binding it.
func adapterInfo(deps *deps.Deps, io *smartio.SmartIO, name string, dep string, help string, module string, origin string) api.AdapterInfo {
	return api.AdapterInfo{
		Name:       name,
		Dep:        dep,
		Help:       help,
		Module:     module,
		Origin:     origin,
		Installed:  io.IsDir(utils.AdapterDir(name)),
		Availables: utils.AvailablesBinding(deps, io, name),
	}
}
