//go:build ignore

// This file is an illustrative sample, not part of the build.
// It shows the same factory pattern on both sides of the sandbox wall:
// inside sandbox/, where the carrier is an api struct, and inside
// adapters/, where the carrier is the adapter struct.
package example_factories

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// ---------------------------------------------------------------
// Inside the sandbox: the carrier is the api struct being filled,
// and the state the closure reads is its propagated Deps.
// ---------------------------------------------------------------

// ExampleLibFunctionFactory returns the closure that fills
// api.Lib.ExampleLibFunction, closed over l, so the injected dependency is
// read at call time through l.Deps.
func ExampleLibFunctionFactory(l *api.Lib) func() int {
	return func() int {
		return l.Deps.ExampleDepFunctionA() + 1
	}
}

// New builds the api.Lib entry point and runs every lib factory over it,
// assigning each return value into its matching field. It is the factory
// aggregate — a field left unassigned here stays nil and panics on first
// call.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.ExampleLibFunction = ExampleLibFunctionFactory(&l)
	return l
}

// ---------------------------------------------------------------
// Outside the sandbox: the carrier is the adapter struct, whose Deps
// field is the contract the factories fill. The state the closure reads
// is the adapter's own configuration.
// ---------------------------------------------------------------

// ExampleAdapter fills deps.Deps with a fixed clock.
type ExampleAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// now is adapter-specific configuration, read by the closure at call time.
	now time.Time
}

// ExampleDepFunctionAFactory returns the closure that fills
// deps.Deps.ExampleDepFunctionA, reading the adapter's fixed clock through a.
func ExampleDepFunctionAFactory(a *ExampleAdapter) func() int {
	return func() int {
		return int(a.now.Unix())
	}
}

// NewAdapter creates a deps.Deps backed by the example adapter. It builds the
// adapter instance, runs every field factory over it and assigns its return
// value into the matching field, and returns the **contract struct** — never
// the concrete adapter type.
func NewAdapter(now time.Time) deps.Deps {
	adapter := &ExampleAdapter{now: now}
	adapter.Deps.ExampleDepFunctionA = ExampleDepFunctionAFactory(adapter)
	return adapter.Deps
}
