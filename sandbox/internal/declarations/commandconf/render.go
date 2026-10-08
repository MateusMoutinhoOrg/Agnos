package commandconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/triggerconf"
)

// Render serializes a CommandConf back to the command.yaml shape, leaving out
// every key that holds its default: `strict` when true, an arg's `type` when
// string, a flag's `keys` when it declared none.
func Render(sandbox *api.Sandbox, conf *CommandConf) string {
	obj := sandbox.Deps.Serializables.CreateObject()

	if conf.HasPriority {
		obj.AddItemToObject("priority", conf.Priority)
	}
	if conf.HasSegments {
		obj.AddItemToObject("segments", conf.Segments)
	}
	if !conf.Strict {
		obj.AddItemToObject("strict", false)
	}
	if len(conf.Args) > 0 {
		obj.AddItemToObject("args", argsArray(sandbox, conf.Args))
	}
	if len(conf.Flags) > 0 {
		obj.AddItemToObject("flags", flagsArray(sandbox, conf.Flags))
	}
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

	return sandbox.Deps.Serializables.SerializeToYaml(obj)
}

// argsArray renders the `args` sequence, in declaration order.
func argsArray(sandbox *api.Sandbox, args []Arg) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, arg := range args {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("id", arg.Id)
		entry.AddItemToObject("start", arg.Start)
		entry.AddItemToObject("end", arg.End)
		if arg.Type != "" && arg.Type != DefaultArgType {
			entry.AddItemToObject("type", arg.Type)
		}
		if arg.Trigger.Exists {
			entry.AddItemToObject("trigger", triggerconf.Render(sandbox, arg.Trigger))
		}
		if arg.Required {
			entry.AddItemToObject("required", true)
		}
		if arg.HasDefault {
			entry.AddItemToObject("default", arg.Default)
		}
		if arg.Description != "" {
			entry.AddItemToObject("description", arg.Description)
		}
		arr.AddItemToArray(entry)
	}
	return arr
}

// flagsArray renders the `flags` sequence, in declaration order.
func flagsArray(sandbox *api.Sandbox, flags []Flag) *serializibles.SerializibleObject {
	arr := sandbox.Deps.Serializables.CreateArray()
	for _, flag := range flags {
		entry := sandbox.Deps.Serializables.CreateObject()
		entry.AddItemToObject("id", flag.Id)
		if flag.HasKeys {
			entry.AddItemToObject("keys", stringArray(sandbox, flag.Keys))
		}
		if flag.Type != "" && flag.Type != DefaultFlagType {
			entry.AddItemToObject("type", flag.Type)
		}
		if flag.Required {
			entry.AddItemToObject("required", true)
		}
		if flag.HasDefault {
			entry.AddItemToObject("default", flag.Default)
		}
		if flag.HasMin {
			entry.AddItemToObject("min", flag.Min)
		}
		if flag.HasMax {
			entry.AddItemToObject("max", flag.Max)
		}
		if len(flag.Enum) > 0 {
			entry.AddItemToObject("enum", stringArray(sandbox, flag.Enum))
		}
		if flag.Pattern != "" {
			entry.AddItemToObject("pattern", flag.Pattern)
		}
		if flag.Trigger.Exists {
			entry.AddItemToObject("trigger", triggerconf.Render(sandbox, flag.Trigger))
		}
		if flag.Description != "" {
			entry.AddItemToObject("description", flag.Description)
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
