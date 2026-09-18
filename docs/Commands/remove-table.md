# `remove-table`

Delete one collection from a database

```bash
agnos remove-table --database <database> [--path <path>] [--quiet] <name>
```

Drops one table from the database's specs.yaml and runs build without compiling: every method spelled after the table goes with it. A table another table still links to is refused — point that link elsewhere first.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the table is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the table to delete |

```bash
agnos remove-table url --database app-database
```

Database System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
