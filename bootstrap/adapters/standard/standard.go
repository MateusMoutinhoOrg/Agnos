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

func New(cacheFilePath string) deps.Deps {
	deps := deps.Deps{
		Println:  PrintlnFactory(),
		CacheLib: CacheLibFactory(cacheFilePath),
	}
	return deps
}
