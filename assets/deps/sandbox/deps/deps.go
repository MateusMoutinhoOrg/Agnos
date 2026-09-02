package deps
{{if .DepsLibs}}
import ({{range .DepsLibs}}
	{{.Name}} "{{$.Module}}/sandbox/deps/{{.Name}}"{{end}}
)
{{end}}
type Deps struct {
{{- range .DepsLibs}}
	{{.Title}} {{.Name}}.Lib
{{- end}}
}
