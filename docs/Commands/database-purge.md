# `database-purge`

Remove the database layer and every declared database

```bash
agnos database-purge [--help] [--path <path>] [--quiet]
```

Removes everything the database asset groups installed, sandbox/internal/databases and docs/Databases whole, and writes database: false. The store contract is left in place: database-init installed it, but a dep, once there, is the project's — 'agnos remove-dep database' is what takes it away.

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
agnos database-purge
agnos database-purge --path ./my-project
```

Database · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
