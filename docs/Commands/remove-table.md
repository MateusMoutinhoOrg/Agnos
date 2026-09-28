# `remove-table`

Delete one collection from a database

```bash
agnos remove-table <Name> --database <database> [--help] [--path <path>] [--quiet]
```

Drops one table from the database's specs.yaml and runs build without compiling: every method spelled after the table goes with it. A table another table still links to is refused — point that link elsewhere first.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the table to delete |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--database` | string, required |  | the database (identifier or package name) the table is declared on | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-table url --database app-database
```

Database System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
