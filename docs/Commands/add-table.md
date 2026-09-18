# `add-table`

Declare one collection of records on a database

```bash
agnos add-table --database <database> [--path <path>] [--quiet] <name>
```

Appends one table to the database's specs.yaml and runs build. A table is born with no fields: what it holds is declared one 'agnos add-table-field' at a time, and every method the table generates is spelled after its name.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the table is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new table (e.g. url) |

```bash
agnos add-table url --database app-database
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
