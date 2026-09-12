package api
{{if .HasDeps}}
import (
	"{{.Module}}/sandbox/deps"
)
{{end}}
// Sandbox is the whole library: one field per contract declared in
// sandbox/api/, each filled by its binder. sandbox.New returns it, and nothing
// callable lives outside of it.
type Sandbox struct {
{{- if .HasDeps}}
	// Deps is every capability the sandbox reaches the outside world
	// through. It rides on the api so that a function handed the Sandbox
	// holds the whole of what it needs, and can call another field of the
	// api besides — which is what makes a field a caller replaced take
	// effect everywhere. It is also the one field that does not cross into
	// a consumer: an installed copy of this contract carries the api, never
	// the wiring behind it.
	Deps *deps.Deps
{{- end}}
{{- range .Constructors }}
	{{ . }} {{ . }}
{{- end }}
}
