package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectBindings returns one entry per declared binding, each carrying
// the adapters its binding.yaml selects, for GenerateBindingNewFiles. The
// order is the declaration's, not the listing of adapters/impls: a binding
// is a selection, and two adapters may implement the same contract, so which
// one binds cannot be read off a directory listing.
func CollectBindings(sandbox *api.Sandbox, io *stagedfs.StagedFS) ([]map[string]any, error) {

	var bindings []map[string]any
	for _, name := range utils.DeclaredBindings(sandbox, io) {

		conf, err := utils.LoadBindingConf(sandbox, io, name)
		if err != nil {
			return nil, err
		}

		var adapters []map[string]string
		for _, adapter := range conf.Adapters {
			adapters = append(adapters, map[string]string{
				"Name":  adapter,
				"Title": utils.DepField(sandbox, adapter),
			})
		}

		bindings = append(bindings, map[string]any{
			"BindingName": name,
			"Adapters":    adapters,
		})
	}

	return bindings, nil
}
