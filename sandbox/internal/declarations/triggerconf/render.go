package triggerconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializibles "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
)

// Render renders one `trigger` object: `values` on a one-of, `value` on every
// other type.
func Render(sandbox *api.Sandbox, trigger Trigger) *serializibles.SerializibleObject {
	entry := sandbox.Deps.Serializables.CreateObject()
	entry.AddItemToObject("type", trigger.Type)
	if trigger.Type == OneOf {
		values := sandbox.Deps.Serializables.CreateArray()
		for _, value := range trigger.Values {
			values.AddItemToArray(value)
		}
		entry.AddItemToObject("values", values)
	} else {
		entry.AddItemToObject("value", trigger.Value)
	}
	if trigger.Negate {
		entry.AddItemToObject("negate", true)
	}
	if trigger.IgnoreCase {
		entry.AddItemToObject("ignore-case", true)
	}
	return entry
}
