package std

import (
	"time"

	keepdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/keepdeps"
	requestdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/requestdeps"
	verbdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/verbdeps"
)

// This package is the sandbox's *copy* of the api the process runtime
// exposes — the same mechanic as verbdeps, keepdeps, embeddeps, iodeps and
// requestdeps, for the same reason: reading the clock, writing to stdout or
// stderr, and constructing the embedded argv/schema/HTTP libraries are all
// OS-bound or third-party effects, so `time`, `fmt`, `os` and the embedded
// libraries may not appear inside the sandbox. The contract is restated
// here, and the adapter — which lives outside the sandbox — is what fills
// it.
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

	// NewVerbLib builds an argv parser bound to the given arguments.
	NewVerbLib func(args []string) verbdeps.Lib

	// NewKeepLib builds a schema-database library rooted at the given path.
	NewKeepLib func(basePath string) keepdeps.Lib

	// NewRequest builds an HTTP request bound to the given url.
	NewRequest func(url string) requestdeps.Request
}
