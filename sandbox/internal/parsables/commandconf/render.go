package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	serializibles "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/serializebles"
)

// Render serializes a CommandConf back to the entries.yaml shape.
func Render(deps *deps.Deps, conf *CommandConf) string {
	obj := deps.Serializebles.CreateObject()

	obj.AddItemToObject("identifiers", stringArray(deps, conf.Identifiers))
	obj.AddItemToObject("category", conf.Category)
	obj.AddItemToObject("help", conf.Help)
	if conf.LongDescription != "" {
		obj.AddItemToObject("long-description", conf.LongDescription)
	}
	if len(conf.Examples) > 0 {
		obj.AddItemToObject("examples", stringArray(deps, conf.Examples))
	}
	if conf.Hidden {
		obj.AddItemToObject("hidden", true)
	}
	if len(conf.Flags) > 0 {
		obj.AddItemToObject("flags", fieldsObject(deps, conf.Flags))
	}
	if len(conf.Args) > 0 {
		obj.AddItemToObject("args", fieldsObject(deps, conf.Args))
	}

	return deps.Serializebles.SerializeToYaml(obj)
}

func fieldsObject(deps *deps.Deps, fields []Field) *serializibles.SerializibleObject {
	obj := deps.Serializebles.CreateObject()
	for _, field := range fields {
		entry := deps.Serializebles.CreateObject()
		if len(field.Identifiers) > 0 {
			entry.AddItemToObject("identifiers", stringArray(deps, field.Identifiers))
		}
		entry.AddItemToObject("description", field.Description)
		if len(field.Examples) > 0 {
			entry.AddItemToObject("examples", stringArray(deps, field.Examples))
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
		obj.AddItemToObject(field.Key, entry)
	}
	return obj
}

func stringArray(deps *deps.Deps, values []string) *serializibles.SerializibleObject {
	arr := deps.Serializebles.CreateArray()
	for _, value := range values {
		arr.AddItemToArray(value)
	}
	return arr
}
