//go:build ignore

// This file is an illustrative sample, not part of the build.
package lib

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// ExampleFunctionDepsCallingFactory returns the closure that fills
// api.Lib.ExampleFunctionDepsCalling, closed over l, so the injected
// dependency is reached through l.Deps at call time.
func ExampleFunctionDepsCallingFactory(l *api.Lib) func(i int) int {
	return func(i int) int {
		depResult := l.Deps.ExampleDepFunctionA()
		add := depResult + 10 + i
		return add
	}
}

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it, assigning each return value into its
// matching function field. Adding a function field to api.Lib means adding
// its factory call and assignment here — an unlisted field stays nil and
// panics on first call.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.ExampleFunctionDepsCalling = ExampleFunctionDepsCallingFactory(&l)
	l.NewExampleObject = NewExampleObjectFactory(&l)
	return l
}
