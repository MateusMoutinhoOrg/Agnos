# {{ .Name }} Index
{{ .Description }}

| Doc | Description |
| --- | --- |
{{- range .Docs }}
| [{{ .Name }}]({{ .Link }}) | {{ .Description }} |
{{- end }}
