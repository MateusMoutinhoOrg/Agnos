# Databases

A **database** is `sandbox/internal/databases/<db>/`, declared by `database.yaml` and generated
whole from it — the same relation `commands/<x>/command.yaml` and `routes/<x>/route.yaml` have
with the `new.go` they render.
{{ if .DatabaseDocs }}
| Database | Package | Keys under |
| --- | --- | --- |
{{- range .DatabaseDocs }}
| [`{{ .Type }}`]({{ .Package }}.md) | `{{ .Package }}` | `{{ .KeyPrefix }}` |
{{- end }}
{{- else }}
No database is declared yet. Run `{{.GeneratorName}} add-database <name>`, then
`{{.GeneratorName}} add-table <table> --database <name>`, and every database lands on this page
on the next build.
{{- end }}

## The files

| File | Written by |
| --- | --- |
| `database.yaml` | `add-database`, `add-table`, `add-table-field` and their editors — never by hand |
| `api.go` | generated: `<T>Record`, `<T>Input`, `<T>Filter`, and the `<Db>` struct of function fields |
| `new.go` | generated: the `databasedeps.Props` and the wiring of each function field |
| `methods.go` | generated: the body of every method |
| `methods_custom.go` | **hand-written** — the only escape; no build reads it or rewrites it |

A database is not a surface of `sandbox/api/`: its methods are typed by table, so there is no
`[]Database` standing where `Cli.Commands` stands. Whoever needs one builds it on the spot with
`<db>.New(sandbox)`, which touches no key — building one is free and creates nothing until the
first record is written.

## database.yaml

```yaml
name: app-database
key-prefix: app
tables:
  - name: url
    fields:
      - name: alias
        type: key
        required: true
      - name: link
        type: string
        required: true
      - name: redirects
        type: integer
```

| Key | What it says |
| --- | --- |
| `name` | the database's own name; the directory is that name with dashes turned into underscores |
| `key-prefix` | `Props.Path`, the key prefix every record is written under |
| `tables[].name` | one collection of records; every method it generates is spelled after it |
| `tables[].fields[].name` | one field of that collection |
| `tables[].fields[].type` | `key`, `string`, `integer`, `number`, `link` or `object` |
| `tables[].fields[].required` | an insert must carry it; ignored on an `object` field |
| `tables[].fields[].target` | the table a `link` points at — required there, empty everywhere else |
| `tables[].fields[].fields` | the fields of a nested `object`, one level deep |

## The methods

| From | Method |
| --- | --- |
| every table | `Add<T>(props <T>Input) (<T>Record, error)` |
| every table | `Find<T>ById(id int64) (<T>Record, bool)` |
| a `key` field | `Find<T>By<Field>(value string) (<T>Record, bool)` |
| every table | `List<T>s(filter <T>Filter) ([]<T>Record, error)` |
| every table | `List<T>sPage(offset int, limit int) ([]<T>Record, error)` — `limit` 0 reads to the end |
| every table | `Count<T>() (int, error)` |
| every plain field | `Set<T><Field>(id int64, value <type>) error` |
| every table | `Remove<T>(id int64) error` |
| a `link` field | `Get<T><Field>(id int64) (<Target>Record, bool)` |
| an `object` field | `Add<T><Sub>(parentId int64, props <T><Sub>Input) (<T><Sub>Record, error)`, `List<T><Sub>s(parentId int64) ([]<T><Sub>Record, error)` |

**A `Find` is born of a `key` field alone.** It is the only field the storage indexes, so it is
the only one a direct lookup reaches. A `string`, `integer` or `number` field is reached through
`List<T>s` and nowhere else: generating a `Find<T>By<Field>` that scans the whole table would
sell a scan with the face of an indexed lookup. `<T>Filter` therefore covers **every** plain
field — text takes `<Field>StartsWith` and `<Field>Equals`, a number takes `<Field>Min` and
`<Field>Max`, and a zero value turns its own filter off.

Three rules hold for every generated method:

- **A search answers `(<T>Record, bool)`, a write answers `error`.** What failed is an error, what
  is absent is a `false` — never one `nil` standing for both.
- **`sandbox *api.Sandbox` comes first**, in every function of `methods.go`.
- **No type assertion without `ok`.** Every stored value is converted in the comma-ok form, so a
  value of the wrong type is an error and never a panic.

## The commands

```bash
{{.GeneratorName}} database-init                                             # install the dep, turn the mechanic on
{{.GeneratorName}} add-database app-database --key-prefix app
{{.GeneratorName}} add-table url --database app-database
{{.GeneratorName}} add-table-field alias --database app-database --table url --type key --required
{{.GeneratorName}} show-database app-database                                # read the declaration back
```

`set-table-field` and the `remove-` half of each pair are the inverses. Those commands are the
only editors of a `database.yaml`, and they re-render the whole file.
