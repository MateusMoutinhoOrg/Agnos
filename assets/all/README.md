{{- render (printf "%s/docs/ReadmeHeader.md" .ConfigDir) }}

## Documentation Index

| Name | Description |
| --- | --- |
{{- range .Themes }}
| [{{ .Name }}](/docs/Index/{{ .Name }}.md) | {{ .Description }} |
{{- end }}

## License


{{ copy "LICENSE" -}}
