package standard

import (
	"fmt"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/bootstrap/sandbox/contracts/deps/agnosdeps"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

// PrintlnFactory returns the closure that fills deps.Deps.Println, writing a
// line to standard output.
func PrintlnFactory() func(a ...any) {
	return func(a ...any) {
		fmt.Println(a...)
	}
}

func CacheLibFactory(cacheFilePath string) agnosdeps.Lib {
	inner := agnoslib.New(agnosadapter.New(cacheFilePath))
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
	deps := deps.Deps{
		Println:  PrintlnFactory(),
		CacheLib: CacheLibFactory(cacheFilePath),
	}
	return deps
}
