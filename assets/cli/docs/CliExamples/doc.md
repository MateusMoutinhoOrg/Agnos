# CliExamples

`examples/cli/<name>/example.sh` is a shell session a reader can copy, and
`{{.Name}} exec-test` runs it and compares what it produced with the golden beside it. What
the two sides share — the `exec-test` flags, the `result.yaml` schema, the normalization of
`cli-output` and the rule that a name present on both sides leaves the same tree — is in
[LibExamples](../LibExamples/doc.md) and is not repeated here.

## Layout

```
examples/cli/<name>/
  example.sh     run with `sh example.sh`, cwd = this directory
  result.yaml    the golden, generated
  TestDir/       the only place the example may write; git-ignored
```

Exactly one `example.sh` per directory, writing nothing outside its own `TestDir`.

## The `{{.Name}}` on the PATH

The script types `{{.Name}}`, because that is what a reader types — an example reads as
documentation, not as a test script. `exec-test` is what makes that name resolve: before the
suite runs it writes an executable named after the `name` of
`{{.ConfigDir}}/project.yaml` into `release/exec-test/`,

```sh
#!/bin/sh
exec go run <project>/cmd/main "$@"
```

and puts that directory in front of the PATH of every run. So an example is always checked
against the code in the tree, never against an installed binary — the same reason `build` is
bootstrapped. The `go run` alias costs no build step of its own; its price is the
`exit status N` line `go run` writes to stderr when the program exits non-zero, which is part
of the `cli-output` of an example that fails on purpose.

The alias carries the project's `name` verbatim, so an example must type exactly that name: in
a project named with a capital, an `example.sh` typing it in lowercase passes on macOS and
fails on Linux.

## Managing examples

```bash
{{.Name}} add-cli-example <name>       # writes example.sh, a stub that already runs
{{.Name}} remove-cli-example <name>    # deletes the directory whole
```

Never create or delete an example directory by hand. `add-cli-example` refuses in a project
with no cli — there would be no binary for the script to type.
{{ if .CliExamples }}
## Declared examples

| example | directory |
|---|---|
{{- range .CliExamples }}
| `{{ . }}` | `examples/cli/{{ . }}/` |
{{- end }}
{{ else }}
No example is declared yet: `examples/cli/` is created by the first
`{{.Name}} add-cli-example`.
{{ end }}
