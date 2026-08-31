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
	Deps deps.Deps

	output io.Writer
}

// NowFactory returns the closure that fills deps.Deps.Now, returning the
// real current time.
func Now() time.Time {
	return time.Now()
}

// PrintfFactory returns the closure that fills deps.Deps.Printf, writing one
// formatted message to the process's standard output. It is what the
// command-line interface inside the sandbox reports through.
func Printf(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// ErrorFactory returns the closure that fills deps.Deps.Error, writing one
// formatted message to the process's standard error.
func Error(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// ErrorfFactory returns the closure that fills deps.Deps.Errorf, formatting
// an error message and returning it as an error.
func Errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

func New() deps.Deps {
	deps := deps.Deps{}
	deps.Now = func() time.Time {
		return Now()
	}
	deps.Printf = func(format string, a ...any) (n int, err error) {
		return Printf(os.Stdout, format, a...)
	}
	deps.Error = func(format string, a ...any) (n int, err error) {
		return Error(os.Stderr, format, a...)
	}
	deps.Errorf = func(format string, a ...any) error {
		return Errorf(format, a...)
	}
	deps.NewVerbLib = verblib.NewVerbLib
	deps.NewKeepLib = keeplib.NewKeepLib
	deps.EmbedDeps = embeddeps.NewEmbedDeps()
	deps.IoLib = iolib.NewIoLib()
	deps.NewRequest = requestlib.NewRequestFactory()
	deps.SerializeLib = serializelib.NewSerializeLib()
	return deps
}
