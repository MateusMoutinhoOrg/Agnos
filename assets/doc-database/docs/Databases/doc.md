# Databases

A **database** is `sandbox/internal/databases/<db>/`, declared by `database.yaml` and generated
whole from it — the same relation `commands/<x>/command.yaml` and `routes/<x>/route.yaml` have
with the `generated.new.go` they render.
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
| `generated.api.go` | generated: `<T>Record`, `<T>Input`, `<T>Filter`, and the `<Db>` struct of function fields |
| `generated.new.go` | generated: the `databasedeps.Props` and the wiring of each function field |
| `generated.methods.go` | generated: the body of every method |
| `generated.migrations.go` | generated: every version, its schema and what brings the data to it — see [Migrations](#migrations) |
| `migrations/v<N>.yaml`, `migrations/v<N>.ops.yaml` | the editing commands: the schema of a past version, the operations of one |
| `migrations/v<N>.go` | **hand-written**, written once: a version's own step |
| `methods_custom.go` | **hand-written** — the only escape; no build reads it or rewrites it |

A database is not a surface of `sandbox/api/`: its methods are typed by table, so there is no
`[]Database` standing where `Cli.Commands` stands. Whoever needs one builds it on the spot with
`<db>.New(sandbox)`, which touches no key — building one is free and creates nothing until the
first record is written.

## Where the data lives

Every database lives in one folder, `data` unless a command line says otherwise:
`<project> start-server --database var/app` writes every database under `./var/app`. The folder is
relative to the directory the program runs from, its segments lower-case letters, digits, `-` and
`_` — the store resolves every key under that directory and spells any other character escaped,
so `--database` refuses anything else. A database's `key-prefix` is a folder inside it
(`data/app` for `key-prefix: app`).

| Piece | What |
| --- | --- |
| `sandbox/api/generated.databaseconfig.go` | generated: `sandbox.Config.DatabaseDir`, and `api.DefaultDatabaseDir` (`data`), the value `NewConfig` starts it at |
| `sandbox/internal/commands/middleware/database_dir/` | generated while `cli` is on: a middleware in front of every command line that reads `--database` into `sandbox.Config.DatabaseDir`; no command declares `--database` itself |

A `key-prefix` under `data/` is a declaration from before `--database`: it would now read
`data/data/<x>`, so `verify` refuses it — write `<x>`.

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
| `key-prefix` | the folder, inside the `--database` folder, every record is written under: lower-case letters, digits, `-` and `_`, never `data/…`, `agnos-migrating/…` or `agnos-pre-migration/…` |
| `version` | the version of the stored data the code expects; `1` when left out, raised by the editing commands alone |
| `tables[].name` | one collection of records; every method it generates is spelled after it |
| `tables[].fields[].name` | one field of that collection |
| `tables[].fields[].type` | `key`, `string`, `integer`, `number`, `bytes`, `link` or `object` |
| `tables[].fields[].required` | an insert must carry it; ignored on an `object` field |
| `tables[].fields[].target` | the table a `link` points at — required there, empty everywhere else |
| `tables[].fields[].fields` | the fields of a nested `object`, one level deep; an `object` is never named `position` or `values` |

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
the only one a direct lookup reaches. A `key` is unique without regard to case: `Ana` and `ana`
are one value. A `string`, `integer` or `number` field is reached through
`List<T>s` and nowhere else: generating a `Find<T>By<Field>` that scans the whole table would
sell a scan with the face of an indexed lookup. `<T>Filter` therefore covers **every** plain
field — text takes `<Field>StartsWith` and `<Field>Equals`, a `bytes` field the same pair as
`[]byte`, a number takes `<Field>Min` and `<Field>Max`, and a zero value turns its own filter off.

**A `bytes` field is a `[]byte`, stored exactly as given** — no text encoding, so a file, an
image or a hash comes back byte for byte. Like a `string` it carries no index.

Three rules hold for every generated method:

- **A search answers `(<T>Record, bool)`, a write answers `error`.** What failed is an error, what
  is absent is a `false` — never one `nil` standing for both.
- **`sandbox *api.Sandbox` comes first**, in every function of `generated.methods.go`.
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

## Migrations

The stored data has a version, recorded in the data itself (the table `agnos-version`, under the
`key-prefix`, so it moves with every copy, backup and restore of the folder). An edit that
changes what stored data holds opens the next version on its own: it freezes `database.yaml` as
it was into `migrations/v<N>.yaml`, writes the operations that carry the data across into
`migrations/v<N+1>.ops.yaml`, and raises `version`. The command says it is the edit, so a rename
is never read as a removal plus an addition.

| Edit | Opens a version | The migration |
| --- | --- | --- |
| `add-table`, `add-table-field` (optional), `--clear required`, a new `target` | no | — a field never written reads as its zero value |
| `add-table-field --required`, `set-table-field --required` | yes | fills every record holding none with `--default <v>` (never on a `key`), or leaves it to the step `--fill-by-hand` writes; a record still holding none fails it |
| `set-table-field --rename` | yes | moves every value to the new name; an `object` moves whole, ids kept |
| `set-table-field --type` into `key`; into `integer` or `link` from any other type; into `number` from `key`, `string` or `bytes` | yes | converts every value; one that does not convert, or two records sharing a key, fails it |
| `set-table-field --type` into `string` or `bytes`; between `integer` and `link`; from either into `number` | no | — every stored value already reads as the new type |
| `set-table-field --type` to or from `object` | yes | deletes what the old shape stored |
| `remove-table-field`, `remove-table` | yes | deletes what the field or the table stored |
| `{{.GeneratorName}} add-database-migration <db>` | yes | nothing but `migrations/v<N>.go`: a seed, a field computed from others |

A version's own step is `V<N>(sandbox *api.Sandbox, from databasedeps.Database, to databasedeps.Database) error`
in `migrations/v<N>.go`, written once and then the project's. It runs after that version's
operations, on the same folder read with the schema of `N-1` (`from`) and of `N` (`to`), and it
names nothing but `databasedeps`, so it compiles whatever the schema becomes.
{{ if .HasCli }}
```bash
{{.ProjectName}} migrate-database              # every database; or name one: migrate-database app-database
```

| Step | What it does |
| --- | --- |
| 1 | copies `<--database>/<key-prefix>` into `<--database>/agnos-migrating/<key-prefix>` |
| 2 | runs every version past the one the data records, on the copy, recording each one reached |
| 3 | repairs every table of the last schema — a field that became a `key` is indexed, and two records holding one value fail — and checks that every record holds every required field |
| 4 | moves the data it replaces to `agnos-pre-migration/<key-prefix>-v<N>` and the copy in its place |

A failure in 1–3 throws the copy away: the data in use was only ever read. A run cut short is
settled by the next one — a copy that passed step 3 is swapped in, any other thrown away.
Migrations go forward only: undoing one is moving the `agnos-pre-migration` folder back.

The `database-version` middleware runs in front of every command line but `migrate-database`,
`help` and `version`, after `--database` is read, and refuses one when a database's data is at an
older version than the code expects, at a newer one — a newer build migrated it — or when a
migration of it stopped half way. Data with no version record is at version 1, written before
versions existed; a database holding no data is stamped with the version it is declared at. It
checks when a command starts, never again: stop a running server, migrate, then start it.
{{ end }}
A project with no cli runs the same engine itself, over `databases.Migrations(sandbox)` from
`sandbox/internal/databases/generated.migrations.go`: `sandbox.Deps.OpinionatedAgnosDatabase.Migrate`
migrates and `CheckVersions` is the guard.
