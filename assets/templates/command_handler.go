package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
)

func CommandHander(deps *deps.Deps, entries *Entries) int {
	deps.Std.Printf("{{.Identifier}} called\n")
	return api.ExitOk
}
