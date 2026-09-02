package standard

import ({{range .Libs}}
	{{.Name}} "{{$.Module}}/adapters/libs/{{.Name}}"{{end}}
	deps "{{.Module}}/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
{{- range .Libs}}
	{{.Name}}.Bind(&deps)
{{- end}}
	return deps
}
