package api

// Cli is the CLI surface of the sandbox. CliMain is the generated
// dispatch-and-parse entry point (see sandbox/internal/cli/climain.go),
// wired in by sandbox/binds/cli.go.
type Cli struct {
	CliMain func(args []string) int
}

const (
	// ExitOk reports that the command did what it was asked to do.
	ExitOk = 0
	// ExitFailure reports that a well-formed command could not be carried
	// out.
	ExitFailure = 1
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command or flag, a missing operand, an unparsable amount. Every such
	// error exits with this one code, whichever command produced it.
	ExitUsage = 2
)
