# `{{ .Command.Identifier }}`
{{- with .Command.Aliases }} — {{ . }}{{ end }}

{{ .Command.Help }}

```bash
{{ .Name }} {{ .Command.Usage }}
```
{{- with .Command.LongDescription }}

{{ . }}
{{- end }}
{{- if .Command.Flags }}

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
{{- range .Command.Flags }}
| {{ .Identifiers }} | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Command.Args }}

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
{{- range .Command.Args }}
| `{{ .Key }}` | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Command.Examples }}

```bash
{{- range .Command.Examples }}
{{ $.Name }} {{ . }}
{{- end }}
```
{{- end }}

{{ .Category }} · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
