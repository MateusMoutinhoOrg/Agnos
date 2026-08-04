package standard

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps/agnosdeps"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

// StandardAdapter fills deps.Deps using the Go standard library plus the
// embedded Agnos cache library, which CacheLibFactory wires up once.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps deps.Deps
	// cacheFilePath is where the embedded cache library persists its records.
	cacheFilePath string
}

// PrintlnFactory returns the closure that fills deps.Deps.Println, writing a
// line to standard output.
func PrintlnFactory(s *StandardAdapter) func(a ...any) {
	return func(a ...any) {
		fmt.Println(a...)
	}
}

// CacheLibFactory returns the value that fills deps.Deps.CacheLib, by
// initializing the embedded Agnos cache library with the embedded library's
// own standard adapter, then copying it onto the agnosdeps.Lib the bootstrap
// sandbox declares. It returns a value instead of a closure, because
// CacheLib is a struct field, not a function field.
//
// Both structs have the same shape, so Set is assigned straight across; only
// Get needs a wrapper, because its Entry is a distinct named type on each side.
// Copying an embedded library into a deps contract is the whole reason the
// contracts are structs of function fields — with interfaces this would need a
// bridging type implementing every method.
//
// Only code outside the sandbox may import the embedded library, so the copy
// happens here.
func CacheLibFactory(s *StandardAdapter) agnosdeps.Lib {
	inner := agnoslib.New(agnosadapter.New(s.cacheFilePath))
	return agnosdeps.Lib{
		Set: inner.Set,
		Get: func(key string) (agnosdeps.Entry, bool) {
			entry, found := inner.Get(key)
			if !found {
				return agnosdeps.Entry{}, false
			}
			return agnosdeps.Entry{
				Value:     entry.Value,
				ExpiresAt: entry.ExpiresAt,
				IsExpired: entry.IsExpired,
			}, true
		},
	}
}

// New creates a deps.Deps backed by the standard adapter, ready for lib.New.
// It builds the adapter instance and runs every field factory over it,
// persisting the embedded library's records as a single JSON file at
// cacheFilePath. Adding a field to deps.Deps means adding its factory call
// here.
func New(cacheFilePath string) deps.Deps {
	adapter := &StandardAdapter{cacheFilePath: cacheFilePath}
	adapter.Deps.Println = PrintlnFactory(adapter)
	adapter.Deps.CacheLib = CacheLibFactory(adapter)
	return adapter.Deps
}
