package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
)

// CommandHandler runs the `{{.Identifier}}` command. Every flag and arg its
// entries.yaml declares is read off command by id — command.GetString("path"),
// command.GetBool("quiet"), command.GetStrings("example").
func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	sandbox.Deps.Std.Printf("{{.Identifier}} called\n")
	return api.ExitOk
}
