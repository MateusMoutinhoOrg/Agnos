# LibExamples

Every example of {{.Name}} used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `TestDir`,
so it can be read as documentation and copied as a starting point.

`{{.Name}} exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one{{ if .HasCli }}; the cli side is [CliExamples](../CliExamples/doc.md){{ end }}.
{{ if .LibExamples }}
| Example | Source |
|---|---|
{{- range .LibExamples }}
| `{{ . }}` | [example.go](../../examples/lib/{{ . }}/example.go) |
{{- end }}
{{ else }}
No example is declared yet: `examples/lib/` is created by the first
`{{.Name}} add-lib-example`.
{{ end }}
