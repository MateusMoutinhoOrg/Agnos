package resultconf

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

func Render(sandbox *api.Sandbox, conf *ResultConf) string {
	obj := sandbox.Deps.SerializableDeps.CreateObject()
	obj.AddItemToObject("cli-output", conf.CliOutput)
	obj.AddItemToObject("exit-code", int64(conf.ExitCode))

	entries := make([]TreeEntry, len(conf.Tree))
	copy(entries, conf.Tree)
	sandbox.Deps.SortDeps.SliceStable(entries, func(i, j int) bool {
		return entries[i].File < entries[j].File
	})

	tree := sandbox.Deps.SerializableDeps.CreateArray()
	for _, entry := range entries {
		item := sandbox.Deps.SerializableDeps.CreateObject()
		item.AddItemToObject("file", entry.File)
		item.AddItemToObject("sha", entry.Sha)
		tree.AddItemToArray(item)
	}
	obj.AddItemToObject("tree", tree)

	return sandbox.Deps.SerializableDeps.SerializeToYaml(obj)
}
