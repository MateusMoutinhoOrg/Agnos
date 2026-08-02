package lib

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/entry"
)

// SetFactory fills api.Lib.Set with a closure that stores value under key,
// expiring ttlSeconds from the injected clock's current time.
func SetFactory(l *api.Lib) func(key string, value string, ttlSeconds int) {
	return func(key string, value string, ttlSeconds int) {
		expiresAt := l.Deps.Now().Add(time.Duration(ttlSeconds) * time.Second)
		l.Deps.Store(key, value, expiresAt.Unix())
	}
}

// GetFactory fills api.Lib.Get with a closure that returns the live Entry
// stored under key. It reports false when the key is absent or when the
// record has already expired.
func GetFactory(l *api.Lib) func(key string) (api.Entry, bool) {
	return func(key string) (api.Entry, bool) {
		value, expiresAtUnix, ok := l.Deps.Load(key)
		if !ok {
			return api.Entry{}, false
		}

		e := entry.New(l.Deps, value, expiresAtUnix)
		if e.IsExpired() {
			return api.Entry{}, false
		}
		return e, true
	}
}

// New builds the api.Lib entry point, storing the injected deps on it and
// running every lib factory over it to fill its function fields. Adding a
// function field to api.Lib means adding its factory call here.
func New(d deps.Deps) api.Lib {
	l := api.Lib{Deps: d}
	l.Set = SetFactory(&l)
	l.Get = GetFactory(&l)
	return l
}
