package version

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
)

// CommandHander backs the `version` / `--version` verb. Entries has no fields:
// this command declares no flags or args.
func CommandHander(deps *deps.Deps, entries *Entries) int {
	if config.Version == "" {
		deps.Std.Printf("no version set yet\n")
		return api.ExitOk
	}
	deps.Std.Printf("Version:%s\n", config.Version)
	return api.ExitOk
}
