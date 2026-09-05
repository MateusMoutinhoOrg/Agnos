package resultconf

import (
	"sort"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

func Render(deps *deps.Deps, conf *ResultConf) string {
	obj := deps.Serializables.CreateObject()
	obj.AddItemToObject("cli-output", conf.CliOutput)
	obj.AddItemToObject("exit-code", int64(conf.ExitCode))

	entries := make([]TreeEntry, len(conf.Tree))
	copy(entries, conf.Tree)
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].File < entries[j].File
	})

	tree := deps.Serializables.CreateArray()
	for _, entry := range entries {
		item := deps.Serializables.CreateObject()
		item.AddItemToObject("file", entry.File)
		item.AddItemToObject("sha", entry.Sha)
		tree.AddItemToArray(item)
	}
	obj.AddItemToObject("tree", tree)

	return deps.Serializables.SerializeToYaml(obj)
}
