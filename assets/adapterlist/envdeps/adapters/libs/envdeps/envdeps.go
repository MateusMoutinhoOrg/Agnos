package envdeps

import (
	"os"

	envdeps "{{.Module}}/sandbox/deps/envdeps"

	"{{.Module}}/sandbox/deps"
)

// Bind fills deps.Deps.Envdeps with the standard library's os.
func Bind(deps *deps.Deps) {
	deps.Envdeps = envdeps.Sandbox{
		Getenv: func(key string) string {
			return os.Getenv(key)
		},
	}
}
