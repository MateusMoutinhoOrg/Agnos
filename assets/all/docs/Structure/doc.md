# Structure

`(gen)` = written by `build`, never edited — the full list is in
[GeneratedFiles](../GeneratedFiles/doc.md).

```
adapters/  -->  sandbox/  <--  cmd/{{ if .HasAssets }}        assets/ (templates, reached via Deps.Embeddeps){{ end }}
(reaches OS)    (closed)       (wires)
```
{{ if .Structure }}
Every line below is one entry of `{{.ConfigDir}}/{{.StructureConfFile}}` — add `<path>:
{description: "..."}` there, nested under `children:` of its parent, with `dir: true` on a
directory, `gen: true` on a file `build` rewrites, and `order:` to place it among its siblings
(unordered siblings follow, alphabetically).

```
{{- range .Structure }}
{{ .Line }}
{{- end }}
```
{{- else }}
Nothing is described yet: add entries to `{{.ConfigDir}}/{{.StructureConfFile}}` and run
`agnos build`.
{{- end }}

Every rule this shape has to hold to — layers, naming, generated files, docs — is in
[Rules](../Rules/doc.md); the command that makes each change is in
[Workflow](../Workflow/doc.md).
