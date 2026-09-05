# CliExamples

Every example of the {{.Name}} cli. Each one is a shell session that runs with its own
directory as the working directory and writes only into its own `TestDir`, so it can be read
as documentation and copied line by line. The script types `{{.Name}}`, which `exec-test`
resolves to the code in this tree.

`{{.Name}} exec-test` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `TestDir` file, written by
`exec-test` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one; the lib side is [LibExamples](../LibExamples/doc.md).
{{ if .CliExamples }}
| Example | Source |
|---|---|
{{- range .CliExamples }}
| `{{ . }}` | [example.sh](../../examples/cli/{{ . }}/example.sh) |
{{- end }}
{{ else }}
No example is declared yet: `examples/cli/` is created by the first
`{{.Name}} add-cli-example`.
{{ end }}
