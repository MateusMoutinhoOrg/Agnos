package std

import (
	"time"
)

// This package is the sandbox's *copy* of the api the process runtime
// exposes — the same mechanic as verbdeps, keepdeps, embeddeps, iodeps and
// requestdeps, for the same reason: reading the clock and writing to stdout
// or stderr are OS-bound effects, so `time`, `fmt` and `os` may not appear
// inside the sandbox. The contract is restated here, and the adapter —
// which lives outside the sandbox — is what fills it.
//
// deps.Deps itself carries no function fields: every capability it exposes
// is a library struct (iodeps.Lib, embeddeps.Lib, …). The loose runtime
// functions the sandbox needs are gathered into this one struct and
// injected whole as the Deps.Std field.

// Lib is the runtime library injected whole as the Deps.Std field.
type Lib struct {
	// Now returns the current wall-clock time.
	Now func() time.Time

	// Printf writes one formatted message to standard output.
	Printf func(format string, a ...any) (n int, err error)

	// Error writes one formatted message to standard error.
	Error func(format string, a ...any) (n int, err error)

	// Errorf formats an error message and returns it as an error.
	Errorf func(format string, a ...any) error
}
