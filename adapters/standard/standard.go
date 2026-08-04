package standard

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps/verbdeps"

	verblib "github.com/MateusMoutinhoOrg/Agnos-Cli/Verb/sandbox"
)

type record struct {
	Value         string `json:"value"`
	ExpiresAtUnix int64  `json:"expiresAtUnix"`
}

// StandardAdapter fills deps.Deps using the Go standard library for storage
// and the clock, and the embedded Verb library — wired over the process's own
// command line — for argument parsing. It persists records in a single JSON
// file configured on New, so values survive across runs, and reads the real
// wall clock for Now. Only files outside the sandbox, like this one, may
// import the embedded Verb library.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps     deps.Deps
	mu       sync.RWMutex
	filePath string
	// args is the argument vector the embedded Verb library parses, taken
	// from the process's own command line.
	args []string
}

// NowFactory returns the closure that fills deps.Deps.Now, returning the
// real current time.
func NowFactory(s *StandardAdapter) func() time.Time {
	return func() time.Time {
		return time.Now()
	}
}

// LoadFactory returns the closure that fills deps.Deps.Load, reading a
// record from the JSON file. ok is false when the file is absent, invalid,
// or the key does not exist.
func LoadFactory(s *StandardAdapter) func(key string) (value string, expiresAtUnix int64, ok bool) {
	return func(key string) (value string, expiresAtUnix int64, ok bool) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		raw, err := os.ReadFile(s.filePath)
		if err != nil {
			return "", 0, false
		}
		var store map[string]record
		if err := json.Unmarshal(raw, &store); err != nil {
			return "", 0, false
		}
		if rec, found := store[key]; found {
			return rec.Value, rec.ExpiresAtUnix, true
		}
		return "", 0, false
	}
}

// StoreFactory returns the closure that fills deps.Deps.Store, writing a
// record to the JSON file.
func StoreFactory(s *StandardAdapter) func(key string, value string, expiresAtUnix int64) {
	return func(key string, value string, expiresAtUnix int64) {
		s.mu.Lock()
		defer s.mu.Unlock()
		store := make(map[string]record)
		if raw, err := os.ReadFile(s.filePath); err == nil {
			_ = json.Unmarshal(raw, &store)
		}
		store[key] = record{Value: value, ExpiresAtUnix: expiresAtUnix}
		if raw, err := json.MarshalIndent(store, "", "  "); err == nil {
			_ = os.WriteFile(s.filePath, raw, 0o644)
		}
	}
}

// VerbLibFactory returns the value that fills deps.Deps.VerbLib: the embedded
// Verb argv-parser library, initialized over the adapter's argument vector,
// copied field by field onto the sandbox's local verbdeps.Lib. It returns a
// value rather than a closure because the field is a struct — see the
// Factories specification.
func VerbLibFactory(s *StandardAdapter) verbdeps.Lib {
	inner := verblib.New(s.args)
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

// New creates a deps.Deps backed by the standard adapter, ready for lib.New.
// Records live as a single JSON file at the provided filePath, and the embedded
// Verb library parses the process's own command line, os.Args[1:]. It builds
// the adapter instance and runs every field factory over it, so each closure
// reads the adapter's state at call time. Adding a field to deps.Deps means
// adding its factory call here.
func New(filePath string) deps.Deps {
	adapter := &StandardAdapter{filePath: filePath, args: os.Args[1:]}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Load = LoadFactory(adapter)
	adapter.Deps.Store = StoreFactory(adapter)
	adapter.Deps.VerbLib = VerbLibFactory(adapter)
	return adapter.Deps
}
