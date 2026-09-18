# Commands
{{ if .CommandDocs }}
`{{.Name}} <command> [flags] [args]`. `{{.Name}} help <command>`, or `{{.Name}} <command> --help`,
prints the same for one command; an empty command line prints the general help and exits 2.
A command declaring a `--help` flag of its own keeps it, and is described through `help` alone.

One page per command, each rendered from that command's `entries.yaml`
([EntriesYaml](../EntriesYaml/doc.md)) on each build — open the one you need rather than this
whole page. Hidden commands are not listed. Flags may appear anywhere on the command line;
positionals bind in order after them. A `repeatable` field is given once per value.
{{- range .CommandDocs }}

## {{ .Category }}

| Command | Does |
| --- | --- |
{{- range .Commands }}
| [`{{ .Identifier }}`]({{ .Identifier }}.md) | {{ .Help }} |
{{- end }}
{{- end }}
{{- else }}
No command is declared yet. Run `{{.GeneratorName}} add-command <name> --help "..." --category "..."`
and every command lands on this page on the next build.
{{- end }}

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
