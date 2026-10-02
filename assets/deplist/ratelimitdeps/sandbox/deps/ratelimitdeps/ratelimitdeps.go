package ratelimitdeps

// This package is the sandbox's *copy* of the api a rate limiter exposes —
// the same mechanic as std and randdeps, for a reason of its own: a limiter
// is state every request of the process shares, written by requests served at
// the same time, and the sandbox may not import `sync` to guard it. The
// contract is restated here, and the adapter — which lives outside the
// sandbox — holds the counters and their lock.
//
// A counter is named by a key the caller chooses ("login-ip:203.0.113.7")
// and counts hits within a fixed window: the first hit of a key opens its
// window, and the first hit after the window closes opens a new one, back
// at one. Counters live in memory, so a restart forgets them all.

// Sandbox is the rate limiter injected whole as the Deps.Ratelimitdeps field.
// Every function is safe to call from requests served concurrently.
type Sandbox struct {
	// Hit records one hit under key, in a window of windowSeconds, and
	// returns how many hits the window holds with it included.
	Hit func(key string, windowSeconds int64) int

	// Count returns how many hits key holds in its open window, 0 when it
	// has none or its window of windowSeconds closed. It records nothing.
	Count func(key string, windowSeconds int64) int

	// Reset forgets every hit of key.
	Reset func(key string)
}
