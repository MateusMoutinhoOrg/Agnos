package hashdeps

import (
	"crypto/sha256"
	"encoding/hex"

	hashdeps "{{.Module}}/sandbox/deps/hashdeps"

	"{{.Module}}/sandbox/deps"
)

// sha256Hex fills hashdeps.Sandbox.Sha256Hex, returning the SHA-256 digest of
// content as lower-case hexadecimal.
func sha256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// Bind fills deps.Deps.Hashdeps with the standard library's crypto/sha256 and
// encoding/hex.
func Bind(deps *deps.Deps) {
	deps.Hashdeps = hashdeps.Sandbox{
		Sha256Hex: func(content []byte) string {
			return sha256Hex(content)
		},
	}
}
