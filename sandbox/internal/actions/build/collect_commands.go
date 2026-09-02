package build

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// CollectCommands returns one {"Name": <dir>, "Title": <Dir>} entry per
// sandbox/internal/commands sub-directory, in listing order, for the
// {{range .Commands}} loop in sandbox/binds/cli.go. Each sub-directory is one
// command package exposing NewCommand(deps, sandbox).
func CollectCommands(io *smartio.SmartIO) []map[string]string {
	return collectLibDirs(io, "sandbox/internal/commands")
}
