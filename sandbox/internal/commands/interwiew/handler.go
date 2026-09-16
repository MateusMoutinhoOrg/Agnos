package interwiew

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// CommandHandler runs the `interwiew` command. Every flag and arg its
// entries.yaml declares is read off command by id — command.GetString("path"),
// command.GetBool("quiet"), command.GetStrings("example").
func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	sandbox.Deps.Std.Printf("interwiew called\n")
	return api.ExitOk
}
