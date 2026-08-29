package keep

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/keepdeps"

	keepadapter "github.com/MateusMoutinhoOrg/Keep/adapters/standard"
	keeplib "github.com/MateusMoutinhoOrg/Keep/sandbox"
	keepapi "github.com/MateusMoutinhoOrg/Keep/sandbox/contracts/api"
)

// toKeepItems converts schema fields from the sandbox's copy into the
// embedded Keep library's own type, recursing into nested (Database)
// fields. The constants of both sides carry the same values, so Type
// crosses over unmapped.
func toKeepItems(items []keepdeps.Item) []keepapi.Item {
	if items == nil {
		return nil
	}
	converted := make([]keepapi.Item, 0, len(items))
	for _, item := range items {
		converted = append(converted, keepapi.Item{
			Name:     item.Name,
			Type:     item.Type,
			Required: item.Required,
			Itens:    toKeepItems(item.Itens),
		})
	}
	return converted
}

// toKeepProps converts a database description from the sandbox's copy into
// the embedded Keep library's own type.
func toKeepProps(props keepdeps.Props) keepapi.Props {
	schemas := make([]keepapi.Schema, 0, len(props.Schemas))
	for _, schema := range props.Schemas {
		schemas = append(schemas, keepapi.Schema{
			Name:  schema.Name,
			Itens: toKeepItems(schema.Itens),
		})
	}
	return keepapi.Props{Path: props.Path, Schemas: schemas}
}

// fromKeepItems converts schema fields the embedded Keep library handed
// back into the sandbox's copy, recursing into nested fields.
func fromKeepItems(items []keepapi.Item) []keepdeps.Item {
	if items == nil {
		return nil
	}
	converted := make([]keepdeps.Item, 0, len(items))
	for _, item := range items {
		converted = append(converted, keepdeps.Item{
			Name:     item.Name,
			Type:     item.Type,
			Required: item.Required,
			Itens:    fromKeepItems(item.Itens),
		})
	}
	return converted
}

// fromKeepError converts a failure the embedded Keep library reported into
// the sandbox's copy. A nil error stays nil — that is how success is
// reported on both sides.
func fromKeepError(err *keepapi.Error) *keepdeps.Error {
	if err == nil {
		return nil
	}
	return &keepdeps.Error{
		Type:     err.Type,
		Key:      err.Key,
		KeyValue: err.KeyValue,
		Message:  err.Message,
	}
}

// fromKeepSchemaItem converts one record the embedded Keep library handed
// back into the sandbox's copy, wrapping every field that returns another
// api struct so nothing of the embedded library reaches the sandbox.
func fromKeepSchemaItem(item keepapi.SchemaItem) keepdeps.SchemaItem {
	return keepdeps.SchemaItem{
		Items:  fromKeepItems(item.Items),
		Prefix: item.Prefix,
		Id:     item.Id,
		Get: func(fieldName string) (any, *keepdeps.Error) {
			value, err := item.Get(fieldName)
			return value, fromKeepError(err)
		},
		Update: func(fieldName string, value any) *keepdeps.Error {
			return fromKeepError(item.Update(fieldName, value))
		},
		Remove: func() *keepdeps.Error {
			return fromKeepError(item.Remove())
		},
		CheckKeysPresence: item.CheckKeysPresence,
		ListAll: func(fieldName string) []keepdeps.SchemaItem {
			return fromKeepSchemaItems(item.ListAll(fieldName))
		},
		NewSubItem: func(fieldName string, fields map[string]any) (keepdeps.SchemaItem, *keepdeps.Error) {
			sub, err := item.NewSubItem(fieldName, fields)
			return fromKeepSchemaItem(sub), fromKeepError(err)
		},
		String: item.String,
	}
}

// fromKeepSchemaItems converts a slice of records into the sandbox's copy.
func fromKeepSchemaItems(items []keepapi.SchemaItem) []keepdeps.SchemaItem {
	if items == nil {
		return nil
	}
	converted := make([]keepdeps.SchemaItem, 0, len(items))
	for _, item := range items {
		converted = append(converted, fromKeepSchemaItem(item))
	}
	return converted
}

// fromKeepSchemaInstance converts one collection the embedded Keep library
// handed back into the sandbox's copy.
func fromKeepSchemaInstance(instance keepapi.SchemaInstance) keepdeps.SchemaInstance {
	return keepdeps.SchemaInstance{
		Items:  fromKeepItems(instance.Items),
		Prefix: instance.Prefix,
		NewItem: func(fields map[string]any) (keepdeps.SchemaItem, *keepdeps.Error) {
			item, err := instance.NewItem(fields)
			return fromKeepSchemaItem(item), fromKeepError(err)
		},
		FindByKey: func(key string, keyValue any) (keepdeps.SchemaItem, bool) {
			item, found := instance.FindByKey(key, keyValue)
			if !found {
				return keepdeps.SchemaItem{}, false
			}
			return fromKeepSchemaItem(item), true
		},
		ListAll: func() ([]keepdeps.SchemaItem, *keepdeps.Error) {
			items, err := instance.ListAll()
			return fromKeepSchemaItems(items), fromKeepError(err)
		},
		List: func(position int, chunk int) ([]keepdeps.SchemaItem, *keepdeps.Error) {
			items, err := instance.List(position, chunk)
			return fromKeepSchemaItems(items), fromKeepError(err)
		},
	}
}

// fromKeepDatabase converts one database the embedded Keep library handed
// back into the sandbox's copy.
func fromKeepDatabase(database keepapi.KeepDatabase) keepdeps.KeepDatabase {
	return keepdeps.KeepDatabase{
		Props: keepdeps.Props{
			Path:    database.Props.Path,
			Schemas: fromKeepSchemas(database.Props.Schemas),
		},
		GetSchema: func(name string) (keepdeps.SchemaInstance, bool) {
			instance, found := database.GetSchema(name)
			if !found {
				return keepdeps.SchemaInstance{}, false
			}
			return fromKeepSchemaInstance(instance), true
		},
	}
}

// fromKeepSchemas converts the collections of a database description into
// the sandbox's copy.
func fromKeepSchemas(schemas []keepapi.Schema) []keepdeps.Schema {
	if schemas == nil {
		return nil
	}
	converted := make([]keepdeps.Schema, 0, len(schemas))
	for _, schema := range schemas {
		converted = append(converted, keepdeps.Schema{
			Name:  schema.Name,
			Itens: fromKeepItems(schema.Itens),
		})
	}
	return converted
}

// NewKeepLib returns the value that fills deps.Deps.KeepLib: the embedded
// Keep schema-database library, wired with Keep's own filesystem adapter over
// the adapter's base directory, copied onto the sandbox's local keepdeps.Lib.
// It returns a value rather than a closure because the field is a struct —
// see the Factories specification.
func NewKeepLib(keepBasePath string) keepdeps.Lib {
	inner := keeplib.New(keepadapter.NewWithBase(keepBasePath))
	return keepdeps.Lib{
		NewDatabase: func(props keepdeps.Props) keepdeps.KeepDatabase {
			return fromKeepDatabase(inner.NewDatabase(toKeepProps(props)))
		},
	}
}
