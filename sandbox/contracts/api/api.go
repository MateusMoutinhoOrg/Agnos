package api

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// Exit codes reported by Lib.Sandboxmain, and through it by the process that
// calls it.
const (
	// ExitOk reports that the requested command ran to completion.
	ExitOk = 0
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command, a missing operand, an unparsable amount — and the usage
	// screen was printed.
	ExitUsage = 1
	// ExitFailure reports that a well-formed command could not be carried
	// out, because a record was missing or could not be written.
	ExitFailure = 2
)

type SandBox struct {
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps

	Sandboxmain func(args []string) int
}
