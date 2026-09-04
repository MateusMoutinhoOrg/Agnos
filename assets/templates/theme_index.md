# {{ .Name }}
{{ .Description }}

| Doc | Description |
| --- | --- |
{{- range .Docs }}
| [{{ .Name }}]({{ .Link }}) | {{ .Description }} |
{{- end }}
