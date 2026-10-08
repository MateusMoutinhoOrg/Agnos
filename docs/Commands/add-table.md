# `add-table`

Declare one collection of records on a database

```bash
agnos add-table <Name> --database <database> [--help] [--path <path>] [--quiet]
```

Appends one table to the database's database.yaml and runs build. A table is born with no fields: what it holds is declared one 'agnos add-table-field' at a time, and every method the table generates is spelled after its name.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new table (e.g. url) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the table is declared on | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-table url --database app-database
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
