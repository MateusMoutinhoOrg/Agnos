//go:build ignore

// This file is an illustrative sample, not part of the build.
package standard

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
)

// StandardAdapter fills deps.Deps using the Go standard library only.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// DepPropA is adapter-specific configuration, read by the closures at
	// call time.
	DepPropA int
}

// ExampleDepFunctionAFactory returns the closure that fills
// deps.Deps.ExampleDepFunctionA, reading the adapter's configuration through
// s, so its state is resolved when the dependency is called.
func ExampleDepFunctionAFactory(s *StandardAdapter) func() int {
	return func() int {
		return s.DepPropA + 20
	}
}

// New creates a deps.Deps backed by the standard adapter.
// depPropA is adapter-specific configuration. It builds the adapter instance,
// runs every field factory over it and assigns its return value into the
// matching field, and returns the contract struct — never the concrete
// adapter type. Adding a field to deps.Deps means adding its factory call
// here.
func New(depPropA int) deps.Deps {
	adapter := &StandardAdapter{DepPropA: depPropA}
	adapter.Deps.ExampleDepFunctionA = ExampleDepFunctionAFactory(adapter)
	return adapter.Deps
}
