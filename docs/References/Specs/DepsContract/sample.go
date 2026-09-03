//go:build ignore

package clockdeps

// This file is an illustrative sample, not part of the build.
//
// This package is the sandbox's *copy* of the api a clock exposes — the same
// mechanic as std, iodeps and rundeps, for the same reason: reading the wall
// clock is an OS-bound effect, so `time.Now` may not appear inside the
// sandbox. The contract is restated here, and the adapter lib — which lives
// outside the sandbox — is what fills it.

import "time"

// Lib is the clock injected whole as the Deps.Clockdeps field.
type Lib struct {
	// Now returns the current wall-clock time.
	Now func() time.Time

	// After returns a Timer that fires once d has elapsed. Timers are
	// created per call, so they are handed back as a struct of function
	// fields rather than injected once.
	After func(d time.Duration) Timer
}

// Timer is one pending fire, handed back by Lib.After.
type Timer struct {
	// Wait blocks until the timer fires or Stop is called, reporting which.
	Wait func() bool
	// Stop cancels the timer. Stopping a fired timer is not an error.
	Stop func()
}