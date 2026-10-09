package sandbox

import (
	api "{{.Module}}/sandbox/api"{{if .HasDeps}}
	deps "{{.Module}}/sandbox/deps"{{end}}
{{- range .ConstructorPackages}}
	{{.}} "{{$.Module}}/sandbox/constructors/{{.}}"
{{- end}}
)

// New builds the whole library: one call per package under
// sandbox/constructors/, each filling the field of the Sandbox it owns. The
// list is the directories themselves, so a constructor written by hand is
// called exactly like a generated one — this file is rendered around what is
// there, never the other way round.
{{if .HasDeps}}func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}
{{else}}func New() *api.Sandbox {
	self := api.Sandbox{}
{{end}}
{{- range .ConstructorPackages}}
	{{.}}.Constructor(&self)
{{- end}}

	return &self
}
