package osstd

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"time"

	stddeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/stddeps"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
)

// now fills stddeps.Contract.Now, returning the real current time as nanoseconds
// since the Unix epoch.
func now() int64 {
	return time.Now().UnixNano()
}

// printf fills stddeps.Contract.Printf, writing one formatted message to the given
// output. It is what the command-line interface inside the sandbox reports
// through.
func printf(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// logWrite fills stddeps.Contract.Log, writing one formatted progress message to the
// given output.
func logWrite(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// errorWrite fills stddeps.Contract.Error, writing one formatted message to the given
// output.
func errorWrite(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// sprintf fills stddeps.Contract.Sprintf, formatting a message and returning it as a
// string.
func sprintf(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

// goos fills stddeps.Contract.GOOS, reporting the operating system the process was
// built for.
func goos() string {
	return runtime.GOOS
}

// errorf fills stddeps.Contract.Errorf, formatting an error message and returning it
// as an error.
func errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// Bind fills deps.Deps.StdDeps with the clock and the three output channels,
// built on the standard library's time and fmt over os.Stdout and os.Stderr.
func Bind(deps *deps.Deps) {
	deps.StdDeps = stddeps.Contract{
		Now: func() int64 {
			return now()
		},
		Printf: func(format string, a ...any) (n int, err error) {
			return printf(os.Stdout, format, a...)
		},
		Logf: func(format string, a ...any) (n int, err error) {
			return logWrite(os.Stderr, format, a...)
		},
		Eprintf: func(format string, a ...any) (n int, err error) {
			return errorWrite(os.Stderr, format, a...)
		},
		Errorf: func(format string, a ...any) error {
			return errorf(format, a...)
		},
		Sprintf: func(format string, a ...any) string {
			return sprintf(format, a...)
		},
		GOOS: func() string {
			return goos()
		},
	}
}
