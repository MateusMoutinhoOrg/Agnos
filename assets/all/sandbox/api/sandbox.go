package api

// Sandbox is the whole library: one field per contract declared in
// sandbox/api/, each filled by its binder. sandbox.New returns it, and nothing
// callable lives outside of it.
type Sandbox struct {
{{- range .Constructors }}
	{{ . }} {{ . }}
{{- end }}
}
