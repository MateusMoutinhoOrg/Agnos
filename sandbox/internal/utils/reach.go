package utils

// MiddlewareReach is whether a middleware runs in front of a command or a
// route, read off the two declarations alone — without a command line or a
// request.
type MiddlewareReach int

const (
	// NoReach is a middleware that never runs in front of the command.
	NoReach MiddlewareReach = iota
	// Runs is a middleware that runs in front of every line of the command.
	Runs
	// MayRun is a middleware whose trigger reads a segment the command
	// captures, or a regex: whether it runs depends on what is typed, and
	// explain-command gives the exact answer.
	MayRun
)
