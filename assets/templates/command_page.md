# `{{ .Command.Identifier }}`
{{- with .Command.Aliases }} — {{ . }}{{ end }}

{{ .Command.Help }}
{{- if .Command.Middleware }}

A middleware: it runs on rung {{ .Command.Priority }}, in front of every command line matching
`{{ .Command.Pattern }}`, and hands the line on unless it answers.
{{- else }}

```bash
{{ .Name }} {{ .Command.Usage }}
```
{{- end }}
{{- with .Command.LongDescription }}

{{ . }}
{{- end }}
{{- if .Command.Args }}

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
{{- range .Command.Args }}
| `{{ .Id }}` | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Command.Flags }}

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
{{- range .Command.Flags }}
| {{ .Keys }} | {{ .Type }} | {{ .Default }} | {{ .Description }} | {{ if .From }}[{{ .From }}]({{ .FromPage }}){{ else }}—{{ end }} |
{{- end }}
{{- end }}
{{- if .Command.Middlewares }}

| Runs in front of it | When |
| --- | --- |
{{- range .Command.Middlewares }}
| [`{{ .Name }}`]({{ .Page }}) | {{ with .Condition }}{{ . }}{{ else }}always{{ end }} |
{{- end }}
{{- end }}
{{- if .Command.Middleware }}

| Runs in front of | When |
| --- | --- |
{{- range .Command.RunsBefore }}
| [`{{ .Name }}`]({{ .Page }}) | {{ with .Condition }}{{ . }}{{ else }}always{{ end }} |
{{- else }}
| — | no command it matches |
{{- end }}
{{- end }}
{{- if .Command.Examples }}

```bash
{{- range .Command.Examples }}
{{ $.Name }} {{ . }}
{{- end }}
```
{{- end }}

{{ .Category }} · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
