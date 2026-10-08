package structureconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

func New(sandbox *api.Sandbox, content string) (*StructureConf, error) {

	if content == "" {
		return nil, sandbox.Deps.Std.Errorf("content cannot be empty, use NewEmpty instead")
	}

	specs, parse_error := sandbox.Deps.Serializables.ParseYaml(content)
	if parse_error != nil {
		return nil, parse_error
	}
	structure_specs := specs

	structure_conf := &StructureConf{
		Items: make([]Item, 0),
	}

	// A file holding nothing but comments parses to null: it declares no item,
	// which is a description of nothing, not a malformed document.
	if structure_specs.IsNull() {
		BindMethods(sandbox, structure_conf)
		return structure_conf, nil
	}

	if !structure_specs.IsObject() {
		return nil, sandbox.Deps.Std.Errorf("structure_specs is not an object")
	}

	items, err := parseItems(sandbox, structure_specs, "")
	if err != nil {
		return nil, err
	}
	structure_conf.Items = items

	BindMethods(sandbox, structure_conf)
	return structure_conf, nil
}

// parseItems reads one object of the document as a set of sibling items,
// recursing into every `children` it finds. parent is the path the items hang
// from, used only to name the element in an error message.
func parseItems(sandbox *api.Sandbox, node *serializibles.SerializibleObject, parent string) ([]Item, error) {
	keys, err := node.GetKeys()
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not get the keys of %s", itemLabel(parent))
	}

	items := make([]Item, 0)

	for _, key := range keys {
		item_specs, _ := node.GetObjectItem(key)
		if item_specs == nil || item_specs.IsNull() {
			continue
		}

		path := key
		if parent != "" {
			path = parent + "/" + key
		}

		if !item_specs.IsObject() {
			return nil, sandbox.Deps.Std.Errorf("%s is not an object", path)
		}

		item, err := parseItem(sandbox, item_specs, key, path)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	sortItems(sandbox, items)
	return items, nil
}

// parseItem reads one item object: its description, its optional gen and order
// keys, and the children nested under it.
func parseItem(sandbox *api.Sandbox, item_specs *serializibles.SerializibleObject, name string, path string) (Item, error) {
	item := Item{
		Name:     name,
		Children: make([]Item, 0),
	}
	var err error

	description_item, _ := item_specs.GetObjectItem("description")
	dir_item, _ := item_specs.GetObjectItem("dir")
	gen_item, _ := item_specs.GetObjectItem("gen")
	order_item, _ := item_specs.GetObjectItem("order")
	children_item, _ := item_specs.GetObjectItem("children")

	if description_item != nil && !description_item.IsNull() {
		item.Description, err = description_item.GetString()
		if err != nil {
			return item, sandbox.Deps.Std.Errorf("%s: description is not a string", path)
		}
	}

	if dir_item != nil && !dir_item.IsNull() {
		item.Dir, err = dir_item.GetBool()
		if err != nil {
			return item, sandbox.Deps.Std.Errorf("%s: dir is not a bool", path)
		}
	}

	if gen_item != nil && !gen_item.IsNull() {
		item.Gen, err = gen_item.GetBool()
		if err != nil {
			return item, sandbox.Deps.Std.Errorf("%s: gen is not a bool", path)
		}
	}

	// order is optional: absent means "list this item after every ordered
	// sibling".
	if order_item != nil && !order_item.IsNull() {
		order, err := order_item.GetInt()
		if err != nil {
			return item, sandbox.Deps.Std.Errorf("%s: order is not an int", path)
		}
		item.Order = int(order)
		item.HasOrder = true
	}

	if children_item != nil && !children_item.IsNull() {
		if !children_item.IsObject() {
			return item, sandbox.Deps.Std.Errorf("%s: children is not an object", path)
		}

		item.Children, err = parseItems(sandbox, children_item, path)
		if err != nil {
			return item, err
		}
	}

	return item, nil
}

// sortItems orders siblings the way the tree renders them: by `order`, then by
// name. An item with no `order` comes after every ordered one.
func sortItems(sandbox *api.Sandbox, items []Item) {
	sandbox.Deps.Sortdeps.SliceStable(items, func(i, j int) bool {
		left, right := items[i], items[j]
		if left.HasOrder != right.HasOrder {
			return left.HasOrder
		}
		if left.HasOrder && left.Order != right.Order {
			return left.Order < right.Order
		}
		return left.Name < right.Name
	})
}

// itemLabel names the element an error is about: the document itself when the
// path is empty.
func itemLabel(path string) string {
	if path == "" {
		return "the structure document"
	}
	return path
}
