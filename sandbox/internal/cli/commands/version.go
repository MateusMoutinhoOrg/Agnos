package commands

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// VersionCommand prints the interface version, read from its own asset so a release
// bump is a one-line edit outside the code. Both the `version` command and
// the --version flag land here.
func VersionCommand(l *api.Lib) int {
	l.Deps.Printf(config.VersionMessage+"\n", Version(l))
	return api.ExitOk
}

// Version returns the interface version the `version` command reports,
// trimmed of the newline its asset file ends with.
func Version(l *api.Lib) string {
	return strings.TrimSpace(config.Version)
}
