package deps
{{if .DepLibs}}
import ({{range .DepLibs}}
	{{.Name}} "{{$.Module}}/sandbox/deps/{{.Name}}"{{end}}
)
{{end}}
// Deps is every capability the sandbox needs from the outside world, one field
// per sub-contract directory of sandbox/deps/. An adapter fills the fields; the
// sandbox only calls them, which is what keeps it free of OS packages.
type Deps struct {
{{- range .DepLibs}}
	{{.Title}} {{.Name}}.{{.Type}}
{{- end}}
}
