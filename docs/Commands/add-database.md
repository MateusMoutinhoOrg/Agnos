# `add-database`

Declare a new database in the project

```bash
agnos add-database [--prefix <prefix>] [--path <path>] [--quiet] <name>
```

Writes sandbox/internal/databases/<name>/specs.yaml and runs build, which generates api.go, new.go and methods.go beside it. A database is born with no tables: 'agnos add-table' declares the first one.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--prefix` | string |  | the key prefix every record is written under (defaults to the database's own name) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new database (e.g. app-database) |

```bash
agnos add-database app-database --prefix app
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
