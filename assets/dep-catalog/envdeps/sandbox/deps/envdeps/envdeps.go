package envdeps

// This package is the sandbox's *copy* of the api the process environment
// exposes — the same mechanic as std and randdeps, for the same reason: the
// sandbox may import nothing but the sandbox, so `os` may not appear inside
// it. The contract is restated here, and the adapter — which lives outside
// the sandbox — is what fills it.
//
// It is how a secret reaches the process without travelling on the command
// line, where every user of the machine can read it, or in a file, which can
// end up committed next to the code.

// Sandbox is the process environment injected whole as the Deps.Envdeps field.
type Sandbox struct {
	// Getenv returns the value of the environment variable key, or "" when
	// it is not set.
	Getenv func(key string) string
}
