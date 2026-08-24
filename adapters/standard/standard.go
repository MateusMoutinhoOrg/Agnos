package standard

import (
	"fmt"
	"io"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// StandardAdapter fills deps.Deps using the Go standard library for the
// clock, the embedded Verb library — wired over the process's own command
// line — for argument parsing, and the embedded Keep library — wired to
// Keep's own filesystem adapter — for the schema database every category and
// transaction is persisted in. Records live on disk under a base directory
// configured on New, so a tracked budget survives across runs. Only files
// outside the sandbox, like this one, may import the embedded Verb and Keep
// libraries.
//
// It also fills the three fields the tracker itself never calls — EmbedDeps
// from the project's compiled-in assets, IoLib from `os` and `path/filepath`,
// and NewRequest from `net/http`. They are capabilities the template offers a
// derived library, and an adapter must fill every field of the contract
// whether the current library exercises it or not: an unfilled field is a nil
// function the compiler will not catch.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// args is the argument vector the embedded Verb library parses, taken
	// from the process's own command line.
	args []string
	// output is the stream deps.Deps.Printf writes to — the process's
	// standard output, which is what a command-line interface reports on.
	output io.Writer
	// keepBasePath is the directory the embedded Keep library writes its
	// records under, one file per key.
	keepBasePath string
}

// NowFactory returns the closure that fills deps.Deps.Now, returning the
// real current time.
func NowFactory(s *StandardAdapter) func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}

// PrintfFactory returns the closure that fills deps.Deps.Printf, writing one
// formatted message to the process's standard output. It is what the
// command-line interface inside the sandbox reports through.
func PrintfFactory(s *StandardAdapter) func(format string, a ...any) (int, error) {
	return func(format string, a ...any) (int, error) {
		return fmt.Fprintf(s.output, format, a...)
	}
}

