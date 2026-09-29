# `add-database`

Declare a new database in the project

```bash
agnos add-database <Name> [--prefix <prefix>] [--help] [--path <path>] [--quiet]
```

Writes sandbox/internal/databases/<name>/specs.yaml and runs build, which generates api.go, new.go and methods.go beside it. A database is born with no tables: 'agnos add-table' declares the first one.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new database (e.g. app-database) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--prefix` | string |  | the key prefix every record is written under (defaults to the database's own name) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-database app-database --prefix app
```

Database System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
