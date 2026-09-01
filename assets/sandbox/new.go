package sandbox

import (
	api "{{.Module}}/sandbox/api"
)

func New() *api.Sandbox {
	self := api.Sandbox{}

	return &self
}
