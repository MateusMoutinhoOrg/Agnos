package sandbox

import (
	api "{{.Module}}/sandbox/api"{{if .HasDeps}}
	deps "{{.Module}}/sandbox/deps"{{end}}
)

{{if .HasDeps}}func New(deps *deps.Deps) *api.Sandbox {
{{else}}func New() *api.Sandbox {
{{end}}	self := api.Sandbox{}

	return &self
}
