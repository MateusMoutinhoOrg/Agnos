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

## Verify rules

- `sandbox/` holds only `api`, `binds`, `deps`, `internal` and `new.go`.
- Nothing under `sandbox/` imports a module package outside `sandbox/`.
- `sandbox/api/*` imports only `sandbox/api`. `sandbox/deps/*` imports only stdlib and `sandbox/deps`.
- Every file of `sandbox/api/` and `sandbox/deps/` parses, and every exported type, func, const and var in them carries a doc comment — `docs/PublicApi/doc.md` is generated from those comments.
- Every `sandbox/binds/` file mirrors an `api/` file and declares only functions.
- `adapters/` holds only `availables` and `libs`. Every `libs/<lib>/` exports `Bind(deps *deps.Deps)`, and every `sandbox/deps/<x>` contract has a lib mentioning its `deps.<X>` field.
- Every `assets/deplist/<dep>/<path>`, rendered with this module, equals `<path>` whenever that file exists here.
- `docs/Commands/doc.md`, when present, names every visible command and every flag it declares.
- Every `docs/**` dir has a parsable `props.yaml`; first-level docs name at least one theme from `themes.yaml`; sub-docs name none. A theme no doc names renders no README section and is not an error.
- Every entry of `{{.ConfigDir}}/{{.StructureConfFile}}` names a path that exists — a directory when it declares `dir: true`, a file otherwise. An entry holding `<`, `*` or `?` stands for a family of paths, and only the literal head of its path is required to exist.

## Naming

- `Deps` field = title-cased `sandbox/deps/<dir>` (`iodeps` -> `deps.Iodeps`). Always use that spelling.
- Adapter lib binder is always `Bind(deps *deps.Deps)`. Command handler is always `CommandHandler(deps *deps.Deps, entries *Entries) int`.
- An adapter lib's file is named after its package (`adapters/libs/iodeps/iodeps.go`), like a contract's (`sandbox/deps/iodeps/iodeps.go`). A second file in the package is named after what it holds.
- A dep is named after the contract it installs, not its lib (`argvdeps`, lib `verb`).
