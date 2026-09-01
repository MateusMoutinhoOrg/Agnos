package api

type Sandbox struct {
{{- range .Constructors }}
	{{ . }} {{ . }}
{{- end }}
}
