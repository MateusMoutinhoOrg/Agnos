package api

import "github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"

// Lib is the entry point handed back by lib.New, a struct of function
// fields filled by the factories in sandbox/internal/lib, each closing over
// this struct so it reads the injected Deps at call time.
type Lib struct {
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps
	// TestFunc exercises the embedded financial-tracker library reached
	// through the deps.
	TestFunc func()
}
