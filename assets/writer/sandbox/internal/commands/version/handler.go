package version

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
	"{{.Module}}/sandbox/internal/config"
)

// CommandHandler backs the `version` / `--version` verb. Entries has no fields:
// this command declares no flags or args.
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	if config.Version == "" {
		deps.Std.Printf("no version set yet\n")
		return api.ExitOk
	}
	deps.Std.Printf("Version:%s\n", config.Version)
	return api.ExitOk
}
