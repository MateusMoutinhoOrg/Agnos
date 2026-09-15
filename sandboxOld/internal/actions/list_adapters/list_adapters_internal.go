package list_adapters

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListAdaptersInternal returns one row per adapter the project can reach: every
// adapter of the embedded catalog, plus every one installed that the catalog
// does not have — a generated shim is only ever in the second group.
func ListAdaptersInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) ([]api.AdapterInfo, error) {
	sandbox.Deps.Std.Log("list-adapters started with path %s \n", path)

	catalog, err := utils.CatalogAdapters(sandbox)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	adapters := []api.AdapterInfo{}

	for _, name := range catalog {
		seen[name] = true
		conf, err := utils.LoadCatalogAdapterConf(sandbox, name)
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, adapterInfo(sandbox, io, name, conf.Dep, conf.Help, conf.Module, conf.Origin))
	}

	for _, name := range utils.InstalledAdapters(sandbox, io) {
		if seen[name] {
			continue
		}
		conf, err := utils.LoadAdapterConf(sandbox, io, name)
		if err != nil {
			continue
		}
		adapters = append(adapters, adapterInfo(sandbox, io, name, conf.Dep, conf.Help, conf.Module, conf.Origin))
	}

	return adapters, nil
}

// adapterInfo fills one row with what the project holds of that adapter: an
// installed package, and the availables binding it.
func adapterInfo(sandbox *api.Sandbox, io *smartio.SmartIO, name string, dep string, help string, module string, origin string) api.AdapterInfo {
	return api.AdapterInfo{
		Name:       name,
		Dep:        dep,
		Help:       help,
		Module:     module,
		Origin:     origin,
		Installed:  io.IsDir(utils.AdapterDir(name)),
		Availables: utils.AvailablesBinding(sandbox, io, name),
	}
}
