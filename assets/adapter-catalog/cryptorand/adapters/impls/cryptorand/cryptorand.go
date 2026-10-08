package randdeps

import (
	"crypto/rand"
	"encoding/hex"

	randdeps "{{.Module}}/sandbox/deps/randdeps"

	"{{.Module}}/sandbox/deps"
)

// randomHex fills randdeps.Sandbox.Hex: bytes bytes read from crypto/rand,
// lower-case hexadecimal.
func randomHex(bytes int) (string, error) {
	buffer := make([]byte, bytes)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buffer), nil
}

// Bind fills deps.Deps.Randdeps with the standard library's crypto/rand and
// encoding/hex.
func Bind(deps *deps.Deps) {
	deps.Randdeps = randdeps.Sandbox{
		Hex: func(bytes int) (string, error) {
			return randomHex(bytes)
		},
	}
}
