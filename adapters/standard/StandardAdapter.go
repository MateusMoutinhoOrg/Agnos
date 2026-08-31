package standard

import (
	"fmt"
	"io"
	"os"
	"time"

	embeddeps "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/embeddeps"
	iolib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/iodeps"
	keeplib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/keepdeps"
	requestlib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/requestdeps"
	serializelib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/serializebles"
	verblib "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard/deps/verbdeps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

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
func Now() time.Time {
	return time.Now()
}

// PrintfFactory returns the closure that fills deps.Deps.Printf, writing one
// formatted message to the process's standard output. It is what the
// command-line interface inside the sandbox reports through.
func Printf(s *StandardAdapter, format string, a ...any) (int, error) {
	return fmt.Fprintf(s.output, format, a...)
}

// ErrorFactory returns the closure that fills deps.Deps.Error, writing one
// formatted message to the process's standard error.
func Error(s *StandardAdapter, format string, a ...any) (int, error) {
	return fmt.Fprintf(os.Stderr, format, a...)
}

// ErrorfFactory returns the closure that fills deps.Deps.Errorf, formatting
// an error message and returning it as an error.
func Errorf(s *StandardAdapter, format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func New(basePath string) deps.Deps {
	adapter := &StandardAdapter{
		args:         os.Args[1:],
		output:       os.Stdout,
		keepBasePath: basePath,
	}
	adapter.Deps.Now = func() time.Time {
		return Now()
	}
	adapter.Deps.Printf = func(format string, a ...any) (n int, err error) {
		return Printf(adapter, format, a...)
	}
	adapter.Deps.Error = func(format string, a ...any) (n int, err error) {
		return Error(adapter, format, a...)
	}
	adapter.Deps.Errorf = func(format string, a ...any) error {
		return Errorf(adapter, format, a...)
	}
	adapter.Deps.VerbLib = verblib.NewVerbLib(adapter.args)
	adapter.Deps.KeepLib = keeplib.NewKeepLib(adapter.keepBasePath)
	adapter.Deps.EmbedDeps = embeddeps.NewEmbedDeps()
	adapter.Deps.IoLib = iolib.NewIoLib()
	adapter.Deps.NewRequest = requestlib.NewRequestFactory()
	adapter.Deps.SerializeLib = serializelib.NewSerializeLib()
	return adapter.Deps
}
