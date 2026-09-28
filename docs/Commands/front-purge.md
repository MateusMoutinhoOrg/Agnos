# `front-purge`

Remove the html front layer from the project

```bash
agnos front-purge [--help] [--path <path>] [--quiet]
```

Removes sandbox/internal/generated/frontio and the frontend route, then rebuilds. assets/frontend/ is left untouched: every file there is the project's content, so front-init puts the route back over it.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos front-purge
```

Front System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
