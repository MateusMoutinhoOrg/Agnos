//go:build ignore

// This file is an illustrative sample, not part of the build.
package api

import "github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"

// ExampleLibObject is an object handed back by the library, created
// through Lib.NewExampleObject with the deps already wired in. Its
// function fields are filled by the factories in
// sandbox/internal/example_object/.
type ExampleLibObject struct {
	// Deps is the dependency set propagated from the lib that created
	// this object, carried here so its factories can reach it.
	Deps deps.Deps
	// FirstProp and SecondProp are plain data fields, set at construction.
	FirstProp  int
	SecondProp string
	// ExampleObjectMethod is filled by example_object.ExampleObjectMethodFactory.
	ExampleObjectMethod func() string
}

// Lib is the entry point handed back by lib.New. Every object it
// creates carries the same deps it was built with.
type Lib struct {
	// Deps is the dependency set injected by lib.New, carried here so
	// every factory-built function field can reach it.
	Deps deps.Deps
	// NewExampleObject is filled by lib.NewExampleObjectFactory.
	NewExampleObject func(firstProp int, secondProp string) ExampleLibObject
	// ExampleFunctionDepsCalling is filled by lib.ExampleFunctionDepsCallingFactory.
	ExampleFunctionDepsCalling func(i int) int
}
