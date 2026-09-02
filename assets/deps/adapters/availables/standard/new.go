package standard

import ({{range .AdapterLibs}}
	{{.Name}} "{{$.Module}}/adapters/libs/{{.Name}}"{{end}}
	deps "{{.Module}}/sandbox/deps"
)

func New() deps.Deps {
	deps := deps.Deps{}
{{- range .AdapterLibs}}
	{{.Name}}.Bind(&deps)
{{- end}}
	return deps
}
