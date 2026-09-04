# Structure

`(gen)` = written by `build`, never edited.

```
adapters/  -->  sandbox/  <--  cmd/        assets/ (templates, reached via Deps.Embeddeps)
(reaches OS)    (closed)       (wires)
```
{{ if .Structure }}
Every line below is declared in `{{.ConfigDir}}/{{.StructureConfFile}}`, one entry per element
worth describing. Add or drop an entry there and run `build`; `verify` rejects an entry whose
path no longer exists, so this tree cannot drift from the disk.

```
{{- range .Structure }}
{{ .Line }}
{{- end }}
```
{{- else }}
Nothing is described yet: add entries to `{{.ConfigDir}}/{{.StructureConfFile}}` and run
`build`.
{{- end }}

Every rule this shape has to hold to — layers, naming, generated files, docs — is in
[Rules](../Rules/doc.md).
