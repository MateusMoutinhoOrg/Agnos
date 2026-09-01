package standard

import (
	deps "{{.Module}}/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
	return deps
}
