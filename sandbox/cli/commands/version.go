package commands

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// VersionCommand prints the interface version, held in its own file under
// sandbox/config so a release bump is a one-line edit touching no logic. Both
// the `version` command and the --version flag land here.
func VersionCommand(l *api.Lib) int {
	l.Deps.Printf(config.VersionMessage+"\n", Version())
	return api.ExitOk
}

// Version returns the interface version the `version` command reports,
// trimmed of any surrounding whitespace the constant is written with.
func Version() string {
	return strings.TrimSpace(config.Version)
}
