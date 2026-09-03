package api

// Cli is the CLI surface of the sandbox. CliMain is the generated
// dispatch-and-parse entry point (see sandbox/internal/cli/climain.go),
// wired in by sandbox/binds/cli.go.
type Cli struct {
	CliMain func(args []string) int
}

const (
	ExitOk = 0
	// ExitUsage reports that the command line itself was wrong — an unknown
	// command, a missing operand, an unparsable amount.
	ExitUsage = 1
	// ExitFailure reports that a well-formed command could not be carried
	// out.
	ExitFailure = 2
)
