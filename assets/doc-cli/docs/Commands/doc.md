# Commands
{{ if .CommandDocs.Groups }}
`{{.Name}} <command> [args] [flags]`. `{{.Name}} help <command>`, or `{{.Name}} <command> --help`,
prints the same for one command; an empty command line prints the general help and exits 2.
A command declaring a `--help` flag of its own keeps it, and is described through `help` alone.

One page per command, each rendered from that command's `command.yaml`
([CommandYaml](../CommandYaml/doc.md)) on each build — open the one you need rather than this
whole page. Hidden commands are not listed. The args are the leading words of the command line;
the flags follow them, in any order. A `repeatable` flag is given once per value. A command's page
lists the flags of the middlewares in front of it too.
{{- range .CommandDocs.Groups }}

## {{ .Category }}

| Command | Does |
| --- | --- |
{{- range .Commands }}
| [`{{ .Identifier }}`]({{ .Page }}) | {{ .Help }} |
{{- end }}
{{- end }}
{{- if .CommandDocs.Middlewares }}

## Middlewares

Run in front of the commands they match, lowest `priority` first; typed by nobody.

| Middleware | Runs before | Priority | Flags it adds |
| --- | --- | --- | --- |
{{- range .CommandDocs.Middlewares }}
| [`{{ .Identifier }}`]({{ .Page }}) | `{{ .Pattern }}` | {{ .Priority }} | {{ range $i, $f := .Flags }}{{ if $i }}, {{ end }}{{ $f.Keys }}{{ else }}—{{ end }} |
{{- end }}
{{- end }}
{{- else }}
No command is declared yet. Run `{{.GeneratorName}} add-command <name> --help "..." --category "..."`
and every command lands on this page on the next build.
{{- end }}

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
