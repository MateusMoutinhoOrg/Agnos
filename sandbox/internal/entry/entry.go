package entry

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

// IsExpiredFactory returns the closure that fills api.Entry.IsExpired, so the
// injected clock in e.Deps is consulted at call time rather than at
// construction time.
func IsExpiredFactory(e *api.Entry) func() bool {
	return func() bool {
		return e.Deps.Now().After(e.ExpiresAt)
	}
}

// New builds an api.Entry from the record's stored fields, propagating the
// library's Deps into it, and runs every entry factory over it to fill its
// function fields. Adding a function field to api.Entry means adding its
// factory call here.
func New(d deps.Deps, value string, expiresAtUnix int64) api.Entry {
	e := api.Entry{
		Deps:      d,
		Value:     value,
		ExpiresAt: time.Unix(expiresAtUnix, 0),
	}
	e.IsExpired = IsExpiredFactory(&e)
	return e
}
