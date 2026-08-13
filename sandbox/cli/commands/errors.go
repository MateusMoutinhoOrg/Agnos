package commands

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
)

// UsageError reports a command line that could not be understood, printing
// the named message — with whatever values it names rendered into it — and
// then the usage screen. The message asset carries the quotes a value is
// shown in, so an empty or space-carrying argument is still visible.
func UsageError(l *api.Lib, format string, a ...any) int {
	l.Deps.Printf("%s ", config.ErrorPrefix)
	l.Deps.Printf(format+"\n\n", a...)
	l.Deps.Printf("%s", config.Usages)
	return api.ExitUsage
}

// Failure reports a well-formed command that could not be carried out,
// printing the named message without the usage screen — the command line was
// fine, the records were not.
func Failure(l *api.Lib, format string, a ...any) int {
	l.Deps.Printf("%s ", config.ErrorPrefix)
	l.Deps.Printf(format+"\n", a...)
	return api.ExitFailure
}
