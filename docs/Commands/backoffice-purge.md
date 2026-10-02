# `backoffice-purge`

Remove the admin backoffice from the project

```bash
agnos backoffice-purge [--help] [--path <path>] [--quiet]
```

Removes every file backoffice-init wrote, with the files the build generated beside them, then rebuilds. The server, front and database layers stay, and so do the deps it installed and the ./backofficedb store holding the users: removing data is left to you.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos backoffice-purge
```

Backoffice System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
