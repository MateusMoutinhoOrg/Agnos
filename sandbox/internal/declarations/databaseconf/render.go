package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// Render serializes a DatabaseConf back to the specs.yaml shape. Keys come out
// in one fixed order and comments are dropped, which is what makes a re-render
// idempotent — and why specs.yaml is never edited by hand.
func Render(sandbox *api.Sandbox, conf *DatabaseConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	obj.AddItemToObject("name", conf.Name)
	obj.AddItemToObject("prefix", conf.Prefix)
	obj.AddItemToObject("tables", tablesArray(sandbox, conf.Tables))

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}

// tablesArray renders `tables` as the ordered sequence it is.
func tablesArray(sandbox *api.Sandbox, tables []Table) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, table := range tables {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("name", table.Name)
		entry.AddItemToObject("fields", fieldsArray(sandbox, table.Fields))
		arr.AddItemToArray(entry)
	}
	return arr
}

// fieldsArray renders one table's — or one nested collection's — fields.
func fieldsArray(sandbox *api.Sandbox, fields []Field) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, field := range fields {
		arr.AddItemToArray(fieldObject(sandbox, field))
	}
	return arr
}

// fieldObject renders one field entry: the keys that say nothing are left out,
// so a plain string field reads as its name and its type alone.
func fieldObject(sandbox *api.Sandbox, field Field) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()
	entry.AddItemToObject("name", field.Name)
	entry.AddItemToObject("type", field.Type)
	if field.Required {
		entry.AddItemToObject("required", true)
	}
	if field.Target != "" {
		entry.AddItemToObject("target", field.Target)
	}
	if field.Type == FieldDatabase {
		entry.AddItemToObject("fields", fieldsArray(sandbox, field.Fields))
	}
	return entry
}
