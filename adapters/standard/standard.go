package standard

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/contracts/deps"
)

type record struct {
	Value         string `json:"value"`
	ExpiresAtUnix int64  `json:"expiresAtUnix"`
}

// StandardAdapter fills deps.Deps using the Go standard library only.
// It persists records in a single JSON file configured on New, so values survive across runs,
// and reads the real wall clock for Now.
type StandardAdapter struct {
	// Deps is the contract this adapter fills; its factories assign into it.
	Deps     deps.Deps
	mu       sync.RWMutex
	filePath string
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

// New creates a deps.Deps backed by the standard adapter, ready for lib.New.
// Records live as a single JSON file at the provided filePath. It builds the
// adapter instance and runs every field factory over it, so each closure reads
// the adapter's state at call time. Adding a field to deps.Deps means adding
// its factory call here.
func New(filePath string) deps.Deps {
	adapter := &StandardAdapter{filePath: filePath}
	adapter.Deps.Now = NowFactory(adapter)
	adapter.Deps.Load = LoadFactory(adapter)
	adapter.Deps.Store = StoreFactory(adapter)
	return adapter.Deps
}
