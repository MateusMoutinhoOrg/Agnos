package deps
{{if .Libs}}
import ({{range .Libs}}
	{{.Name}} "{{$.Module}}/sandbox/deps/{{.Name}}"{{end}}
)
{{end}}
type Deps struct {
{{- range .Libs}}
	{{.Title}} {{.Name}}.Lib
{{- end}}
}
