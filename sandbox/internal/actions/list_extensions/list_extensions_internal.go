package list_extensions

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// ListExtensionsInternal returns one row per mechanic of the catalog, in the
// order the catalog declares them, saying which ones this project turned on.
func ListExtensionsInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) ([]api.ExtensionInfo, error) {
	extensions_conf, err := utils.LoadExtensionsConf(sandbox, io)
	if err != nil {
		return nil, err
	}

	utils.NormalizeExtensions(extensions_conf)

	specs := utils.ExtensionCatalog()
	rows := make([]api.ExtensionInfo, 0, len(specs))

	for _, spec := range specs {
		rows = append(rows, api.ExtensionInfo{
			Name:    spec.Name,
			Enabled: extensions_conf.IsEnabled(spec.Name),
			Help:    spec.Help,
		})
	}

	return rows, nil
}
