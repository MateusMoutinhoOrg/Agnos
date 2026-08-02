package deps

import "time"

// Deps is the dependency contract every adapter must satisfy. It is a
// struct of function fields, not an interface: an adapter fills every
// field with the behavior it provides, and the library calls those fields
// directly.
//
// Each field is one injectable behavior the library needs. For this TTL
// cache the injectable behaviors are the clock (so expiry can be tested
// without real waiting) and the storage backend (so records can live in
// memory, on disk, or anywhere an adapter decides).
type Deps struct {
	// Now returns the current time, used to stamp and check expiry.
	Now func() time.Time
	// Load fetches a stored record by key. ok is false when the key is absent.
	Load func(key string) (value string, expiresAtUnix int64, ok bool)
	// Store persists a record and the unix timestamp at which it expires.
	Store func(key string, value string, expiresAtUnix int64)
}
