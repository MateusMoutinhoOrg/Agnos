package agnosdeps

import "time"

// This package is the bootstrap library's *copy* of the embedded Agnos
// cache library's public api. The sandbox may not import the embedded
// library — that would be a third-party import — so it restates the shape
// it needs here, field for field. The adapter, which lives outside the
// sandbox, is what fills these structs from the real library.
//
// Copying is cheap precisely because the embedded library exposes structs
// of function fields instead of interfaces: an adapter assigns the real
// library's fields straight into the copy.

// Entry mirrors the embedded library's api.Entry.
type Entry struct {
	// Value is the cached value.
	Value string
	// ExpiresAt is the moment the entry stops being valid.
	ExpiresAt time.Time
	// IsExpired reports whether the entry has passed ExpiresAt.
	IsExpired func() bool
}

// Lib mirrors the embedded library's api.Lib — a key/value cache with
// per-key expiry.
type Lib struct {
	// Set stores value under key, expiring ttlSeconds from now.
	Set func(key string, value string, ttlSeconds int)
	// Get returns the live Entry stored under key. The bool is false when
	// the key is absent or when the record has already expired.
	Get func(key string) (Entry, bool)
}
