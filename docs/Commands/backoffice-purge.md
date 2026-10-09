# `backoffice-purge`

Remove the admin backoffice from the project

```bash
agnos backoffice-purge [--help] [--path <path>] [--quiet]
```

Removes every file backoffice-init wrote, with the files the build generated beside them, then rebuilds. The server, front and database layers stay, and so do the deps it installed and the ./data/backofficedb and ./data/backup stores holding the users and their backups: removing data is left to you.

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
agnos backoffice-purge
```

Backoffice · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
