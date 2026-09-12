package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
)

func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	sandbox.Deps.Std.Printf("{{.Identifier}} called\n")
	return api.ExitOk
}
