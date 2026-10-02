package timedeps

import (
	"time"

	timedeps "{{.Module}}/sandbox/deps/timedeps"

	"{{.Module}}/sandbox/deps"
)

// formatUnix fills timedeps.Sandbox.FormatUnix: the instant seconds, in UTC,
// spelled by layout.
func formatUnix(seconds int64, layout string) string {
	return time.Unix(seconds, 0).UTC().Format(layout)
}

// parseUnix fills timedeps.Sandbox.ParseUnix: value read as a UTC date spelled
// by layout, in seconds since the Unix epoch.
func parseUnix(layout string, value string) (int64, error) {
	parsed, err := time.ParseInLocation(layout, value, time.UTC)
	if err != nil {
		return 0, err
	}
	return parsed.Unix(), nil
}

// Bind fills deps.Deps.Timedeps with the standard library's time.
func Bind(deps *deps.Deps) {
	deps.Timedeps = timedeps.Sandbox{
		FormatUnix: func(seconds int64, layout string) string {
			return formatUnix(seconds, layout)
		},
		ParseUnix: func(layout string, value string) (int64, error) {
			return parseUnix(layout, value)
		},
	}
}
