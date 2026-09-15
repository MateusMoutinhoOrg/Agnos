package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// Render serializes a CommandConf back to the entries.yaml shape.
func Render(sandbox *api.Sandbox, conf *CommandConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	obj.AddItemToObject("identifiers", stringArray(sandbox, conf.Identifiers))
	obj.AddItemToObject("category", conf.Category)
	obj.AddItemToObject("help", conf.Help)
	if conf.LongDescription != "" {
		obj.AddItemToObject("long-description", conf.LongDescription)
	}
	if len(conf.Examples) > 0 {
		obj.AddItemToObject("examples", stringArray(sandbox, conf.Examples))
	}
	if conf.Hidden {
		obj.AddItemToObject("hidden", true)
	}
	if len(conf.Flags) > 0 {
		obj.AddItemToObject("flags", fieldsArray(sandbox, conf.Flags))
	}
	if len(conf.Args) > 0 {
		obj.AddItemToObject("args", fieldsArray(sandbox, conf.Args))
	}

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}

// fieldsArray renders flags/args in the canonical ordered sequence shape.
func fieldsArray(sandbox *api.Sandbox, fields []Field) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, field := range fields {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("name", field.Key)
		if len(field.Identifiers) > 0 {
			entry.AddItemToObject("identifiers", stringArray(sandbox, field.Identifiers))
		}
		if field.Description != "" {
			entry.AddItemToObject("description", field.Description)
		}
		if len(field.Examples) > 0 {
			entry.AddItemToObject("examples", stringArray(sandbox, field.Examples))
		}
		entry.AddItemToObject("type", field.Type)
		if field.HasDefault {
			entry.AddItemToObject("default", field.Default)
		}
		if field.Required {
			entry.AddItemToObject("required", true)
		}
		if field.Array {
			entry.AddItemToObject("array", true)
		}
		if field.HasMin {
			entry.AddItemToObject("min", field.Min)
		}
		if field.HasMax {
			entry.AddItemToObject("max", field.Max)
		}
		arr.AddItemToArray(entry)
	}
	return arr
}

func stringArray(sandbox *api.Sandbox, values []string) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, value := range values {
		arr.AddItemToArray(value)
	}
	return arr
}
