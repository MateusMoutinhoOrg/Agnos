package main

import (
	"fmt"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
	agnostypes "github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/api"
)

// A real service caches the short-lived bearer token it gets from an auth
// server and only re-authenticates once the token expires. That is exactly
// what a TTL cache is for.

const tokenKey = "auth_token"

// authenticate simulates an expensive round-trip to an auth server that
// mints a fresh bearer token. In real code this would be an HTTP call.
func authenticate() string {
	fmt.Println("  → authenticating against the auth server...")
	return fmt.Sprintf("bearer-%d", time.Now().UnixNano())
}

// token returns a valid bearer token, authenticating only when the cached
// one is missing or has expired (the classic cache-aside pattern).
func token(l agnostypes.Lib) string {
	if entry, ok := l.Get(tokenKey); ok {
		return entry.Value // cache HIT — reuse the live token
	}

	fresh := authenticate()
	l.Set(tokenKey, fresh, 2) // this demo's token lives for 2 seconds
	return fresh
}

func main() {
	// Build deps via an adapter (JSON file store + real clock) and inject.
	l := agnoslib.New(agnosadapter.New("authtoken.json"))

	// Three back-to-back requests: only the first authenticates; the
	// other two reuse the cached token.
	for i := 1; i <= 3; i++ {
		fmt.Printf("request %d uses %s\n", i, token(l))
	}

	// Wait for the token to expire, then call again — it re-authenticates.
	fmt.Println("\nwaiting for the token to expire...")
	time.Sleep(3 * time.Second)
	fmt.Printf("request 4 uses %s\n", token(l))
}
