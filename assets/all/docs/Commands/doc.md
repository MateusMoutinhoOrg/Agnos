# Commands

`{{.Name}} <command> [flags] [args]`. `{{.Name}} help <command>` prints the same for one
command; an empty command line prints the general help and exits 2.
{{ if .CommandDocs }}
Every section below is rendered from `sandbox/internal/commands/<name>/entries.yaml` on each
build: declare a command with `add-command`, its fields with `add-flag` / `add-arg`, and its
prose with `set-command`. Hidden commands are not listed. Flags may appear anywhere on the
command line; positionals bind in order after them. A `repeatable` field is given once per
value.
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
No command is declared yet: run `{{.Name}} cli-init`, then `{{.Name}} add-command <name>`.
{{- end }}

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
