package version

import (
	"{{.Module}}/sandbox/api"
)

// CommandHandler backs the `version` / `--version` verb. Nothing is read off
// the command: it declares no flags or args.
func CommandHandler(sandbox *api.Sandbox, command *api.Command) int {
	if sandbox.Config.Version == "" {
		sandbox.Deps.Std.Printf("no version set yet\n")
		return api.ExitOk
	}
	sandbox.Deps.Std.Printf("Version: %s\n", sandbox.Config.Version)
	return api.ExitOk
}
