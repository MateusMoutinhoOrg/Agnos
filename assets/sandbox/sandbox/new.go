package sandbox

import (
	api "{{.Module}}/sandbox/api"{{if .HasDeps}}
	deps "{{.Module}}/sandbox/deps"{{end}}
{{- range .Constructors}}
{{- if .HasNew}}
	{{.Package}} "{{$.Module}}/sandbox/internal/{{.Package}}"
{{- end}}
{{- end}}
)

{{if .HasDeps}}func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}
{{else}}func New() *api.Sandbox {
	self := api.Sandbox{}
{{end}}
{{- range .Constructors}}
{{- if .HasNew}}
	self.{{.Name}} = {{.Package}}.New{{.Name}}(&self)
{{- end}}
{{- end}}

	return &self
}
