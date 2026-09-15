package hashdeps

// This package is the sandbox's *copy* of the api a hashing library exposes —
// the same mechanic as argvdeps, embeddeps, iodeps, rundeps, std, stringsdeps,
// sortdeps and templatedeps, for the same reason: the sandbox may import
// nothing but the sandbox, so `crypto/sha256` and `encoding/hex` may not
// appear inside it. The contract is restated here, and the adapter — which
// lives outside the sandbox — is what fills it.

// Sandbox is the hashing library injected whole as the Deps.Hashdeps field.
type Sandbox struct {
	// Sha256Hex returns the SHA-256 digest of content, lower-case
	// hexadecimal. It is what every recorded example tree is compared by, so
	// the encoding is part of the golden and may not change.
	Sha256Hex func(content []byte) string
}
