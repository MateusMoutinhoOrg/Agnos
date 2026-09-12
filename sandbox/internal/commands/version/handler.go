package version

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
)

// CommandHandler backs the `version` / `--version` verb. Entries has no fields:
// this command declares no flags or args.
func CommandHandler(sandbox *api.Sandbox, entries *Entries) int {
	if config.Version == "" {
		sandbox.Deps.Std.Printf("no version set yet\n")
		return api.ExitOk
	}
	sandbox.Deps.Std.Printf("Version:%s\n", config.Version)
	return api.ExitOk
}
