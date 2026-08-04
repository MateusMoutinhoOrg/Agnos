package memory

import (
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"

	verblib "github.com/MateusMoutinhoOrg/Agnos-Cli/Verb/sandbox"
)

// record is one stored value together with the moment it expires.
type record struct {
	value         string
	expiresAtUnix int64
}

// MemoryAdapter fills deps.Deps keeping everything in memory.
// Records live in a map guarded by a mutex — nothing is persisted, so
// the store vanishes when the process exits — and Now reads the real
// wall clock. It touches no OS state at all: the argument vector the
// embedded Verb library parses is the one the caller hands to New, not the
// process's command line.
type MemoryAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps  deps.Deps
	mu    sync.RWMutex
	store map[string]record
	// args is the argument vector the embedded Verb library parses, provided
	// by the caller instead of read from the process.
	args []string
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

// VerbLibFactory returns the value that fills deps.Deps.VerbLib: the embedded
// Verb argv-parser library, initialized over the argument vector the caller
// passed to New, copied field by field onto the sandbox's local verbdeps.Lib.
// It returns a value rather than a closure because the field is a struct —
// see the Factories specification.
func VerbLibFactory(m *MemoryAdapter) verbdeps.Lib {
	inner := verblib.New(m.args)
	return verbdeps.Lib{
		Args: inner.Args,
		Used: inner.Used,

		IsPresent: inner.IsPresent,

		GetOptionsSize:   inner.GetOptionsSize,
		GetKeyValuesSize: inner.GetKeyValuesSize,

		GetStringOption:    inner.GetStringOption,
		GetIntOption:       inner.GetIntOption,
		GetDoubleOption:    inner.GetDoubleOption,
		GetTimestampOption: inner.GetTimestampOption,

		GetStringArg:    inner.GetStringArg,
		GetIntArg:       inner.GetIntArg,
		GetDoubleArg:    inner.GetDoubleArg,
		GetTimestampArg: inner.GetTimestampArg,

		GetNextStringArg:    inner.GetNextStringArg,
		GetNextIntArg:       inner.GetNextIntArg,
		GetNextDoubleArg:    inner.GetNextDoubleArg,
		GetNextTimestampArg: inner.GetNextTimestampArg,

		GetStringKeyValues:    inner.GetStringKeyValues,
		GetIntKeyValues:       inner.GetIntKeyValues,
		GetDoubleKeyValues:    inner.GetDoubleKeyValues,
		GetTimestampKeyValues: inner.GetTimestampKeyValues,
	}
}

// New creates a deps.Deps backed by the memory adapter, ready for lib.New.
// The embedded Verb library parses the provided args — pass os.Args[1:] from
// outside the adapter to parse the real command line, or a fixed slice in a
// test. It builds the adapter instance and runs every field factory over it,
// so each closure reads the adapter's in-memory map at call time. Adding a
// field to deps.Deps means adding its factory call here.
func New(args []string) deps.Deps {
	adapter := &MemoryAdapter{store: make(map[string]record), args: args}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Load = LoadFactory(adapter)
	adapter.Deps.Store = StoreFactory(adapter)
	adapter.Deps.VerbLib = VerbLibFactory(adapter)
	return adapter.Deps
}
