package sandbox

import (
	api "{{.Module}}/sandbox/api"{{if .HasDeps}}
	deps "{{.Module}}/sandbox/deps"{{end}}{{if .Binds}}

	binds "{{.Module}}/sandbox/binds"{{end}}
)

{{if .HasDeps}}func New(deps *deps.Deps) *api.Sandbox {
{{else}}func New() *api.Sandbox {
{{end}}	self := api.Sandbox{}
{{range .Binds}}	binds.{{.}}({{if $.HasDeps}}deps, {{end}}&self)
{{end}}
	return &self
}
