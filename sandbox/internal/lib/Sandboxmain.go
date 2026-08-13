package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/cli"
)

// SandboxmainFactory fills api.Lib.Sandboxmain with a closure that runs the
// whole command-line interface: it dispatches one command line — read
// through the injected Verb parser, reported through the injected Printf —
// over the library functions below, and returns the process exit code. The
// dispatch itself lives in sandbox/internal/cli, so this file keeps one
// factory per api.Lib field.
func SandboxmainFactory(l *api.Lib) func(args []string) int {
	return func(args []string) int {
		return cli.Run(l, args)
	}
}
