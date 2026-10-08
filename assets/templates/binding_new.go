package {{.Name}}

import ({{range .Adapters}}
	{{.Name}} "{{$.Module}}/adapters/libs/{{.Name}}"{{end}}
	deps "{{.Module}}/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
{{- range .Adapters}}
	{{.Name}}.Bind(&deps)
{{- end}}
	return deps
}
