package databaseconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializabledeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializabledeps"
)

// Render serializes a DatabaseConf back to the database.yaml shape. Keys come out
// in one fixed order and comments are dropped, which is what makes a re-render
// idempotent — and why database.yaml is never edited by hand.
func Render(sandbox *api.Sandbox, conf *DatabaseConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()

	obj.AddItemToObject("name", conf.Name)
	obj.AddItemToObject("key-prefix", conf.KeyPrefix)
	obj.AddItemToObject("tables", tablesArray(sandbox, conf.Tables))

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}

// tablesArray renders `tables` as the ordered sequence it is.
func tablesArray(sandbox *api.Sandbox, tables []Table) *serializabledeps.SerializableObject {
	arr := sandbox.Deps.SerializableDeps.CreateArray()
	for _, table := range tables {
		entry := sandbox.Deps.SerializableDeps.CreateObject()
		entry.AddItemToObject("name", table.Name)
		entry.AddItemToObject("fields", fieldsArray(sandbox, table.Fields))
		arr.AddItemToArray(entry)
	}
	return arr
}

// fieldsArray renders one table's — or one nested collection's — fields.
func fieldsArray(sandbox *api.Sandbox, fields []Field) *serializabledeps.SerializableObject {
	arr := sandbox.Deps.SerializableDeps.CreateArray()
	for _, field := range fields {
		arr.AddItemToArray(fieldObject(sandbox, field))
	}
	return arr
}

// fieldObject renders one field entry: the keys that say nothing are left out,
// so a plain string field reads as its name and its type alone.
func fieldObject(sandbox *api.Sandbox, field Field) *serializabledeps.SerializableObject {
	entry := sandbox.Deps.SerializableDeps.CreateObject()
	entry.AddItemToObject("name", field.Name)
	entry.AddItemToObject("type", field.Type)
	if field.Required {
		entry.AddItemToObject("required", true)
	}
	if field.Target != "" {
		entry.AddItemToObject("target", field.Target)
	}
	if field.Type == FieldObject {
		entry.AddItemToObject("fields", fieldsArray(sandbox, field.Fields))
	}
	return entry
}
