# Commands
{{ if .CommandDocs }}
`{{.Name}} <command> [flags] [args]`. `{{.Name}} help <command>` prints
the same for one command; an empty command line prints the general help and exits 2.

Hidden commands are not listed. Flags may appear anywhere on the command line; positionals
bind in order after them. A `repeatable` field is given once per value. Every section below is
rendered from that command's `entries.yaml` ([EntriesYaml](../EntriesYaml/doc.md)) on each
build.
{{- range .CommandDocs }}

## {{ .Category }}
{{- range .Commands }}

### `{{ .Identifier }}`
{{- with .Aliases }} — {{ . }}{{ end }}

{{ .Help }}

```bash
{{ $.Name }} {{ .Usage }}
```
{{- with .LongDescription }}

{{ . }}
{{- end }}
{{- if .Flags }}

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
{{- range .Flags }}
| {{ .Identifiers }} | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Args }}

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
{{- range .Args }}
| `{{ .Key }}` | {{ .Type }} | {{ .Default }} | {{ .Description }} |
{{- end }}
{{- end }}
{{- if .Examples }}

```bash
{{- range .Examples }}
{{ $.Name }} {{ . }}
{{- end }}
```
{{- end }}
{{- end }}
{{- end }}
{{- else }}
No command is declared yet — this project has no CLI surface. Run `agnos cli-init`, then
`agnos add-command <name> --help "..." --category "..."`, and every command lands on this page
on the next build.
{{- end }}

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
