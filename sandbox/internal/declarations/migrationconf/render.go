package migrationconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// Render serializes a MigrationConf back to the v<N>.ops.yaml shape. Keys come
// out in one fixed order and the ones that say nothing are left out, which is
// what makes a re-render idempotent.
func Render(sandbox *api.Sandbox, conf *MigrationConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()

	ops := sandbox.Deps.SerializableDeps.CreateArray()
	for _, op := range conf.Ops {
		entry := sandbox.Deps.SerializableDeps.CreateObject()
		entry.AddItemToObject("kind", op.Kind)
		entry.AddItemToObject("table", op.Table)
		for _, pair := range [][2]string{
			{"parent", op.Parent},
			{"field", op.Field},
			{"type", op.Type},
			{"to", op.To},
			{"default", op.Default},
		} {
			if pair[1] != "" {
				entry.AddItemToObject(pair[0], pair[1])
			}
		}
		ops.AddItemToArray(entry)
	}
	obj.AddItemToObject("ops", ops)

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
