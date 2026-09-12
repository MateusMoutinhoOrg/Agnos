package structureconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func BindMethods(sandbox *api.Sandbox, structure_conf *StructureConf) {

	structure_conf.GetItem = func(path string) (*Item, error) {
		item := findItem(structure_conf.Items, "", path)
		if item == nil {
			return nil, sandbox.Deps.Std.Errorf("item not found")
		}
		return item, nil
	}

	structure_conf.Render = func() string {
		return Render(sandbox, structure_conf)
	}
}

// findItem walks the tree for the item whose full path — every ancestor's name
// joined by "/" — is the one asked for. A name may itself hold slashes, so the
// path is compared whole rather than segment by segment.
func findItem(items []Item, parent string, path string) *Item {
	for index, item := range items {
		item_path := item.Name
		if parent != "" {
			item_path = parent + "/" + item.Name
		}

		if item_path == path {
			return &items[index]
		}

		if found := findItem(item.Children, item_path, path); found != nil {
			return found
		}
	}
	return nil
}
