package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// 1. Build deps via an adapter (JSON file store + real clock).
	deps := agnosadapter.New("kvcache.json")

	// 2. Inject deps into the pure library — a key/value cache with TTL.
	l := agnoslib.New(deps)

	// 3. Store a value that stays valid for 60 seconds.
	l.Set("greeting", "hello world", 60)

	// 4. Read it back. The library never knows which adapter is behind it.
	if entry, ok := l.Get("greeting"); ok {
		fmt.Printf("value=%q expires=%s\n", entry.Value, entry.ExpiresAt)
	}

	// 5. A key that was never set is a miss.
	if _, ok := l.Get("missing"); !ok {
		fmt.Println("missing: cache miss")
	}
}
