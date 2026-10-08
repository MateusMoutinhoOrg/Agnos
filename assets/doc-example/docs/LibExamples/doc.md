# LibExamples

Every example of {{.ProjectName}} used as a Go module. Each one is a `package main` program that
runs with its own directory as the working directory and writes only into its own `test-dir`,
so it can be read as documentation and copied as a starting point. It ends by copying out of
`test-dir` into `assert-dir` the paths it asserts — `os.CopyFS(dst, os.DirFS(src))`, one call per
path, each keeping the place it holds in the tree.

`{{.GeneratorName}} run-examples` runs them all and checks each against the `result.yaml` beside it — the
golden holding the output, the exit code and the sha256 of every `assert-dir` file, written by
`run-examples` and never by hand. [Workflow](../Workflow/doc.md) has the commands that add and
remove one{{ if .HasCli }}; the cli side is [CliExamples](../CliExamples/doc.md){{ end }}.
{{ if .LibExamples }}
| Example | Description | Source |
|---|---|---|
{{- range .LibExamples }}
| `{{ .Name }}` | {{ .Description }} | [example.go](../../examples/lib/{{ .Name }}/example.go) |
{{- end }}
{{ else }}
No example is declared yet: `examples/lib/` is created by the first
`{{.GeneratorName}} add-lib-example`.
{{ end }}
