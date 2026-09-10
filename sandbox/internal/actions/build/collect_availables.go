package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectAvailables returns one entry per declared available, each carrying
// the adapters its available.yaml selects, for GenerateAvailableNews. The
// order is the declaration's, not the listing of adapters/libs: an available
// is a selection, and two adapters may implement the same contract, so which
// one binds cannot be read off a directory listing.
func CollectAvailables(deps *deps.Deps, io *smartio.SmartIO) ([]map[string]any, error) {

	var availables []map[string]any
	for _, name := range utils.DeclaredAvailables(deps, io) {

		conf, err := utils.LoadAvailableConf(deps, io, name)
		if err != nil {
			return nil, err
		}

		var adapters []map[string]string
		for _, adapter := range conf.Adapters {
			adapters = append(adapters, map[string]string{
				"Name":  adapter,
				"Title": utils.DepField(deps, adapter),
			})
		}

		availables = append(availables, map[string]any{
			"Name":     name,
			"Adapters": adapters,
		})
	}

	return availables, nil
}
