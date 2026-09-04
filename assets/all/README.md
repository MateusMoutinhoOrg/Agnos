{{- render (printf "%s/docs/ReadmeHeader.md" .ConfigDir) }}

## Documentation Index

| Name | Description |
| --- | --- |
{{- range .Themes }}
| [{{ .Name }}](docs/Index/{{ .Id }}.md) | {{ .Description }} |
{{- end }}

## License


{{ copy "LICENSE" -}}
