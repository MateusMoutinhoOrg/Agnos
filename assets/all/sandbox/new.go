package sandbox

import (
	api "{{.Module}}/sandbox/api"{{if .HasDeps}}
	deps "{{.Module}}/sandbox/deps"{{end}}{{if .Binds}}

	binds "{{.Module}}/sandbox/binds"{{end}}
)

{{if .HasDeps}}func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}
{{else}}func New() *api.Sandbox {
	self := api.Sandbox{}
{{end}}{{range .Binds}}	binds.{{.}}(&self)
{{end}}
	return &self
}
