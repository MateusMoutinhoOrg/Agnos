package api

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
)

// Entry is a single cached record handed back by the library. Deps is the
// injected dependency set the entry was built with; plain data fields carry
// the record itself; IsExpired is a function field filled by a factory in
// sandbox/internal/entry, closing over this struct so it reads Deps and
// ExpiresAt at call time.
type Entry struct {
	// Deps is the dependency set injected into the library, carried here so
	// the entry's factories can reach it.
	Deps deps.Deps
	// Value is the cached value.
	Value string
	// ExpiresAt is the moment the entry stops being valid.
	ExpiresAt time.Time
	// IsExpired reports whether the injected clock has passed ExpiresAt.
	IsExpired func() bool
}

// Lib is the entry point handed back by lib.New. It is a key/value cache
// with per-key expiry, exposed as a struct of function fields: lib.New
// stores the injected deps in Deps and then runs the factories in
// sandbox/internal/lib, each of which fills one function field with a
// closure over this struct.
//
// Because it is a struct and not an interface, a consumer that itself uses
// this pattern can copy the shape of Lib into its own deps contract and
// receive the whole library as a single injected field.
type Lib struct {
	// Deps is the dependency set injected by lib.New, carried here so every
	// factory-built function field can reach it.
	Deps deps.Deps
	// Set stores value under key, expiring ttlSeconds from now.
	Set func(key string, value string, ttlSeconds int)
	// Get returns the live Entry stored under key. The bool is false when
	// the key is absent or when the record has already expired.
	Get func(key string) (Entry, bool)
}
