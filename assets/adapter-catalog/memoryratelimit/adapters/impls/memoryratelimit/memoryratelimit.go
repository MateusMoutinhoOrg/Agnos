package memoryratelimit

import (
	"sync"
	"time"

	ratelimitdeps "{{.Module}}/sandbox/deps/ratelimitdeps"

	"{{.Module}}/sandbox/deps"
)

// sweepAbove is how many keys the limiter holds before a hit sweeps out the
// ones whose window closed, so keys a client invents never pile up.
const sweepAbove = 4096

// maxKeys is the most keys the limiter ever holds. A new key past it, once
// the closed windows are swept out, evicts the window opened longest ago:
// memory stays bounded however many keys clients invent, at the cost of that
// one counter starting over.
const maxKeys = 65536

// window is the counter of one key.
type window struct {
	// opened is when the window's first hit came.
	opened time.Time
	// length is how long the window lasts.
	length time.Duration
	// hits is how many hits came within it.
	hits int
}

// limiter holds every key's window behind one lock.
type limiter struct {
	lock    sync.Mutex
	windows map[string]*window
}

// open answers the window of key that is open at now, nil when there is none.
func (limits *limiter) open(key string, now time.Time) *window {
	found, ok := limits.windows[key]
	if !ok || now.Sub(found.opened) >= found.length {
		return nil
	}
	return found
}

// hit fills ratelimitdeps.Contract.Hit.
func (limits *limiter) hit(key string, windowSeconds int64) int {
	limits.lock.Lock()
	defer limits.lock.Unlock()

	now := time.Now()
	found := limits.open(key, now)
	if found == nil {
		if len(limits.windows) >= sweepAbove {
			limits.sweep(now)
		}
		if len(limits.windows) >= maxKeys {
			limits.evictOldest()
		}
		found = &window{opened: now, length: time.Duration(windowSeconds) * time.Second}
		limits.windows[key] = found
	}
	found.hits++
	return found.hits
}

// count fills ratelimitdeps.Contract.Count.
func (limits *limiter) count(key string, windowSeconds int64) int {
	limits.lock.Lock()
	defer limits.lock.Unlock()

	found := limits.open(key, time.Now())
	if found == nil {
		return 0
	}
	return found.hits
}

// reset fills ratelimitdeps.Contract.Reset.
func (limits *limiter) reset(key string) {
	limits.lock.Lock()
	defer limits.lock.Unlock()

	delete(limits.windows, key)
}

// undo fills ratelimitdeps.Contract.Undo.
func (limits *limiter) undo(key string) {
	limits.lock.Lock()
	defer limits.lock.Unlock()

	found := limits.open(key, time.Now())
	if found != nil && found.hits > 0 {
		found.hits--
	}
}

// evictOldest drops the window opened longest ago. The lock is already held.
func (limits *limiter) evictOldest() {
	oldestKey := ""
	var oldest *window
	for key, found := range limits.windows {
		if oldest == nil || found.opened.Before(oldest.opened) {
			oldestKey, oldest = key, found
		}
	}
	if oldest != nil {
		delete(limits.windows, oldestKey)
	}
}

// sweep drops every window closed by now. The lock is already held.
func (limits *limiter) sweep(now time.Time) {
	for key, found := range limits.windows {
		if now.Sub(found.opened) >= found.length {
			delete(limits.windows, key)
		}
	}
}

// Bind fills deps.Deps.RatelimitDeps with fixed-window counters held in
// memory, at most maxKeys of them, guarded by the standard library's
// sync.Mutex.
func Bind(deps *deps.Deps) {
	shared := &limiter{windows: map[string]*window{}}
	deps.RatelimitDeps = ratelimitdeps.Contract{
		Hit: func(key string, windowSeconds int64) int {
			return shared.hit(key, windowSeconds)
		},
		Count: func(key string, windowSeconds int64) int {
			return shared.count(key, windowSeconds)
		},
		Reset: func(key string) {
			shared.reset(key)
		},
		Undo: func(key string) {
			shared.undo(key)
		},
	}
}
