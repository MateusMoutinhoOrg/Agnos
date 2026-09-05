# LibExamples

An example is documentation and a test at once: `examples/lib/<name>/example.go` is a program
a reader can copy, and `{{.Name}} exec-test` runs it and compares what it produced with the
golden beside it.{{ if .HasCli }} The cli side is [CliExamples](../CliExamples/doc.md); the
schema, the normalization and the cross-check below cover both.{{ end }}

## Layout

```
examples/lib/<name>/
  example.go     package main, run with `go run example.go`, cwd = this directory
  result.yaml    the golden, generated
  TestDir/       the only place the example may write; git-ignored
```

Exactly one `example.go` per directory. It writes nothing outside its own `TestDir`, and it
reports failure by panicking — the lib's counterpart of a cli error message.

## Managing examples

```bash
{{.Name}} add-lib-example <name>       # writes example.go, a stub that already runs
{{.Name}} remove-lib-example <name>    # deletes the directory whole
```

Never create or delete an example directory by hand, the same rule `add-doc` / `remove-doc`
hold: both commands run `build` afterwards, so the table below is rewritten from the tree.

## Running the suite

```bash
{{.Name}} exec-test                    # every example, cli side before lib side
{{.Name}} exec-test --only <name>      # one name, both sides
{{.Name}} exec-test --update           # rewrite every golden with what this run produced
```

Per example, in alphabetical order: `TestDir` is removed, the example runs with its own
directory as the working directory, and the result is compared against `result.yaml` — or
written there when the file does not exist yet. Without `--update` the only other way to
refresh a golden is to delete it. Exit `0` when every check passes, `1` when any diverges;
a divergence prints the expected and produced `exit-code`, a line diff of `cli-output`, and
per file of the tree `+` (only produced), `-` (only golden) or `~` (different sha).

## result.yaml

```yaml
cli-output: "start started with path TestDir \n"
exit-code: 0
tree:
    - file: AgnosConfig/project.yaml
      sha: "af0796675f5c16a2f2b5a3285d21406a30875803f12d51885c9cc33a882e6a51"
```

| field | what it is |
|---|---|
| `cli-output` | stdout and stderr, merged in the order they were written, normalized |
| `exit-code` | the status the example exited with; `0` is success |
| `tree` | every file inside `TestDir`, ordered by `file`, with the sha256 of its content |

All three are compared. The file is generated — keys alphabetical, no trailing space — so it
has to come out byte-for-byte identical on a re-run or the suite stops being idempotent.

**Normalization.** Before comparing or writing, the absolute path of the example's own
directory becomes `<dir>` and `\r\n` becomes `\n`. An example whose output carries any other
absolute path, a timestamp or a resolved version is not a valid example.

**Volatility.** `go.sum` and `release/` are left out of the tree: an example that reaches the
Go runtime (`start` runs `build`, which runs `go mod tidy` and `go build`) writes whatever the
module proxy resolved that day, which is not a property of the project.

## cli and lib are the same run

When `<name>` exists on both sides, the two runs must agree on `tree` and on `exit-code` —
that is the assertion that the cli is only a wrapper over the lib. `cli-output` is left out of
that comparison: each side has its own text (a cli error on one, a `panic` on the other) and is
checked against its own golden.
{{ if .LibExamples }}
## Declared examples

| example | directory |
|---|---|
{{- range .LibExamples }}
| `{{ . }}` | `examples/lib/{{ . }}/` |
{{- end }}
{{ else }}
No example is declared yet: `examples/` is created by the first `{{.Name}} add-lib-example`.
{{ end }}
