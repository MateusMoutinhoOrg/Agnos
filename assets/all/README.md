{{- render (printf "%s/docs/ReadmeHeader.md" .ConfigDir) }}

## Documentation
{{ range .DocIndex }}
### {{ .Name }}

{{ .Description }}

| Doc | Description |
| --- | --- |
{{- range .Docs }}
| [{{ .Name }}]({{ .Link }}) | {{ .Description }} |
{{- end }}
{{ end }}
## License

[MIT](LICENSE)
