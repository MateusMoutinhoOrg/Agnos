package signaldeps

// This package is the sandbox's *copy* of the api a process-signal library
// exposes — the same mechanic as serverdeps, iodeps and std, for the same
// reason: hearing a signal is an OS-bound effect, so `os/signal` may not
// appear inside the sandbox. The contract is restated here, and the adapter —
// which lives outside the sandbox — is what fills it.
//
// It hears a signal and nothing more: what the process does about one — stop
// a server, flush a file — is the caller's business.

// Sandbox is the signal library injected whole as the Deps.Signaldeps field.
type Sandbox struct {
	// OnInterrupt calls handler once, on a goroutine of its own, the first
	// time the process is asked to stop — an interrupt (Ctrl+C) or a
	// termination request. It returns at once. A second request to stop,
	// while handler is still running, ends the process the way it would
	// have ended without OnInterrupt.
	OnInterrupt func(handler func())
}
