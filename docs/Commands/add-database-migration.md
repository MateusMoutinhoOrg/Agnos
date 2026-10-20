# `add-database-migration`

Open the next version of a database for a hand-written step

```bash
agnos add-database-migration <Database> [--help] [--path <path>] [--quiet]
```

Opens the next version of the database's database.yaml with no operation of its own — the schema stays, the data changes — and writes its hand-written step, migrations/v<N>.go, once, then runs build. It is the version no editing command opens: a seed, a field computed from others, a value split in two. The project's migrate-database runs the step on a copy of the data, after the steps before it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Database` | string, required |  | the database (identifier or package name) to open the next version of |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-database-migration app-database
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
