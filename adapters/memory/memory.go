package memory

import (
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/contracts/deps"
)

// record is one stored value together with the moment it expires.
type record struct {
	value         string
	expiresAtUnix int64
}

// MemoryAdapter fills deps.Deps keeping everything in memory.
// Records live in a map guarded by a mutex — nothing is persisted, so
// the store vanishes when the process exits — and Now reads the real
// wall clock.
type MemoryAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps  deps.Deps
	mu    sync.RWMutex
	store map[string]record
}

// NowFactory returns the closure that fills deps.Deps.Now, returning the
// real current time.
func NowFactory(m *MemoryAdapter) func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}

// LoadFactory returns the closure that fills deps.Deps.Load, fetching a
// record from the in-memory map.
func LoadFactory(m *MemoryAdapter) func(key string) (value string, expiresAtUnix int64, ok bool) {
	return func(key string) (value string, expiresAtUnix int64, ok bool) {
		m.mu.RLock()
		defer m.mu.RUnlock()
		r, ok := m.store[key]
		return r.value, r.expiresAtUnix, ok
	}
}

// StoreFactory returns the closure that fills deps.Deps.Store, writing a
// record into the in-memory map.
func StoreFactory(m *MemoryAdapter) func(key string, value string, expiresAtUnix int64) {
	return func(key string, value string, expiresAtUnix int64) {
		m.mu.Lock()
		defer m.mu.Unlock()
		m.store[key] = record{value: value, expiresAtUnix: expiresAtUnix}
	}
}

// New creates a deps.Deps backed by the memory adapter, ready for lib.New.
// It builds the adapter instance and runs every field factory over it, so each
// closure reads the adapter's in-memory map at call time. Adding a field to
// deps.Deps means adding its factory call here.
func New() deps.Deps {
	adapter := &MemoryAdapter{store: make(map[string]record)}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Load = LoadFactory(adapter)
	adapter.Deps.Store = StoreFactory(adapter)
	return adapter.Deps
}
