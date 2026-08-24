package standard

import (
	"fmt"
	"io"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"

	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
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

// VerbLibFactory returns the value that fills deps.Deps.VerbLib: the embedded
// Verb argv-parser library, initialized over the adapter's argument vector,
// copied field by field onto the sandbox's local verbdeps.Lib. It returns a
// value rather than a closure because the field is a struct — see the
// Factories specification.
func VerbLibFactory(s *StandardAdapter) verbdeps.Lib {
	inner := verblib.New(s.args)
	return verbdeps.Lib{
		Args: inner.Args,
		Used: inner.Used,

		IsPresent: inner.IsPresent,

		GetOptionsSize:   inner.GetOptionsSize,
		GetKeyValuesSize: inner.GetKeyValuesSize,

		GetStringOption:    inner.GetStringOption,
		GetIntOption:       inner.GetIntOption,
		GetDoubleOption:    inner.GetDoubleOption,
		GetTimestampOption: inner.GetTimestampOption,

		GetStringArg:    inner.GetStringArg,
		GetIntArg:       inner.GetIntArg,
		GetDoubleArg:    inner.GetDoubleArg,
		GetTimestampArg: inner.GetTimestampArg,

		GetNextStringArg:    inner.GetNextStringArg,
		GetNextIntArg:       inner.GetNextIntArg,
		GetNextDoubleArg:    inner.GetNextDoubleArg,
		GetNextTimestampArg: inner.GetNextTimestampArg,

		GetStringKeyValues:    inner.GetStringKeyValues,
		GetIntKeyValues:       inner.GetIntKeyValues,
		GetDoubleKeyValues:    inner.GetDoubleKeyValues,
		GetTimestampKeyValues: inner.GetTimestampKeyValues,
	}
}


