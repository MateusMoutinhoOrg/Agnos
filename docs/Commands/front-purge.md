# `front-purge`

Remove the html front layer from the project

```bash
agnos front-purge [--help] [--path <path>] [--quiet]
```

Removes the front route, then rebuilds. assets/front/ is left untouched: every file there is the project's content, so front-init puts the route back over it. The OpinionatedAgnosFront lib stays installed, like every dep.

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
agnos front-purge
```

Front · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
