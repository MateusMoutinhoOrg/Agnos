package standard

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps/agnosdeps"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

// StandardAdapter fills the bootstrap deps.Deps using the Go standard library
// for output and the embedded Agnos cache library — wired with its own
// standard adapter — for storage. Only this file, which lives outside the
// bootstrap sandbox, may import the embedded library.
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

// CacheLibFactory returns the value that fills deps.Deps.CacheLib: the
// embedded cache library, initialized with its own standard adapter over the
// adapter's cache file, copied field by field onto the sandbox's local
// agnosdeps.Lib. It returns a value rather than a closure because the field
// is a struct — see the Factories specification.
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

// New creates a deps.Deps backed by the bootstrap standard adapter, ready for
// lib.New. The embedded cache library persists its records at the provided
// cacheFilePath. It builds the adapter instance and runs every field factory
// over it, assigning each return value into the matching field. Adding a field
// to deps.Deps means adding its factory call here.
func New(cacheFilePath string) deps.Deps {
	adapter := &StandardAdapter{cacheFilePath: cacheFilePath}
	adapter.Deps.Println = PrintlnFactory(adapter)
	adapter.Deps.CacheLib = CacheLibFactory(adapter)
	return adapter.Deps
}
