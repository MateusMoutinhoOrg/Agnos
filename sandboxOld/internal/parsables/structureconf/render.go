package structureconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

func Render(sandbox *api.Sandbox, structure_conf *StructureConf) string {
	return sandbox.Deps.Serializables.SerializeToYaml(renderItems(sandbox, structure_conf.Items))
}

// renderItems writes one set of siblings back as an object keyed by item name,
// recursing into the children of every item that has some. Keys are written in
// the order the items are held, which New has already sorted.
func renderItems(sandbox *api.Sandbox, items []Item) *serializibles.SerializibleObject {
	items_obj := sandbox.Deps.Serializables.CreateObject()

	for _, item := range items {
		item_obj := sandbox.Deps.Serializables.CreateObject()
		item_obj.AddItemToObject("description", item.Description)

		if item.Dir {
			item_obj.AddItemToObject("dir", item.Dir)
		}

		if item.Gen {
			item_obj.AddItemToObject("gen", item.Gen)
		}

		if item.HasOrder {
			item_obj.AddItemToObject("order", item.Order)
		}

		if len(item.Children) > 0 {
			item_obj.AddItemToObject("children", renderItems(sandbox, item.Children))
		}

		items_obj.AddItemToObject(item.Name, item_obj)
	}

	return items_obj
}
