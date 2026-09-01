package std

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/std"
)

// now fills std.Lib.Now, returning the real current time.
func now() time.Time {
	return time.Now()
}

// printf fills std.Lib.Printf, writing one formatted message to the given
// output. It is what the command-line interface inside the sandbox reports
// through.
func printf(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// errorWrite fills std.Lib.Error, writing one formatted message to the given
// output.
func errorWrite(output io.Writer, format string, a ...any) (int, error) {
	return fmt.Fprintf(output, format, a...)
}

// errorf fills std.Lib.Errorf, formatting an error message and returning it
// as an error.
func errorf(format string, a ...any) error {
	return fmt.Errorf(format, a...)
}

// NewStdLib returns the value that fills deps.Deps.Std: the loose runtime
// functions the sandbox needs, implemented over the standard library's time,
// fmt and os packages. Printf writes to os.Stdout and Error to os.Stderr.
func NewStdLib() std.Lib {
	return std.Lib{
		Now: func() time.Time {
			return now()
		},
		Printf: func(format string, a ...any) (n int, err error) {
			return printf(os.Stdout, format, a...)
		},
		Error: func(format string, a ...any) (n int, err error) {
			return errorWrite(os.Stderr, format, a...)
		},
		Errorf: func(format string, a ...any) error {
			return errorf(format, a...)
		},
	}
}
