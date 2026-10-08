package timedeps

// This package is the sandbox's *copy* of the api a calendar library exposes —
// the same mechanic as std, which hands out the clock, for the same reason:
// the sandbox may not name a `time.Time`, so turning an instant into a date
// and a date back into an instant is restated here, and the adapter — which
// lives outside the sandbox — is what fills it.
//
// An instant crosses this boundary as seconds since the Unix epoch, and a
// layout is spelled the way Go's `time` package spells one ("2006-01-02").
// Every date is read and written in UTC.

// Sandbox is the calendar library injected whole as the Deps.Timedeps field.
type Sandbox struct {
	// FormatUnix returns the instant seconds, in UTC, spelled by layout.
	FormatUnix func(seconds int64, layout string) string

	// ParseUnix reads value as a UTC date spelled by layout and returns it in
	// seconds since the Unix epoch. The error reports a value that does not
	// fit layout.
	ParseUnix func(layout string, value string) (int64, error)
}
