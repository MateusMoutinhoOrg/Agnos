package clockdeps

// This file is an illustrative sample, not part of the build.

import (
	"time"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	clockdeps "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/clockdeps"
)

// Bind fills deps.Deps.Clockdeps with the real wall clock over the standard
// library's time package. Timers are created per call by after.
func Bind(deps *deps.Deps) {
	deps.Clockdeps = clockdeps.Lib{
		Now:   time.Now,
		After: after,
	}
}

// after fills clockdeps.Lib.After, wrapping one time.Timer onto the sandbox's
// Timer copy. Nothing of the standard library's type crosses the wall.
func after(d time.Duration) clockdeps.Timer {
	inner := time.NewTimer(d)
	stopped := make(chan struct{})
	return clockdeps.Timer{
		Wait: func() bool {
			select {
			case <-inner.C:
				return true
			case <-stopped:
				return false
			}
		},
		Stop: func() {
			if inner.Stop() {
				close(stopped)
			}
		},
	}
}
