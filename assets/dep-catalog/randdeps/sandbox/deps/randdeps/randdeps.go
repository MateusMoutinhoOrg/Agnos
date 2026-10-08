package randdeps

// This package is the sandbox's *copy* of the api a cryptographic random
// source exposes — the same mechanic as hashdeps and jwtdeps, for the same
// reason: the sandbox may import nothing but the sandbox, so `crypto/rand`
// may not appear inside it. The contract is restated here, and the adapter —
// which lives outside the sandbox — is what fills it.

// Sandbox is the random source injected whole as the Deps.Randdeps field.
// Every byte it hands out comes from the operating system's cryptographically
// secure generator, so what it answers may be used as a secret.
type Sandbox struct {
	// Hex returns bytes random bytes, lower-case hexadecimal — a string twice
	// bytes long. The error reports a generator that could not be read.
	Hex func(bytes int) (string, error)
}
