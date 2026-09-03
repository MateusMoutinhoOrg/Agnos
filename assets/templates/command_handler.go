package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
)

func CommandHandler(deps *deps.Deps, entries *Entries) int {
	deps.Std.Printf("{{.Identifier}} called\n")
	return api.ExitOk
}
