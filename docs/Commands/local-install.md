# `local-install`

Builds the project and installs it locally

```bash
agnos local-install [--help] [--path <path>] [--quiet]
```

Runs build over the project, then compiles ./cmd/main into /usr/local/bin/<project-name> (~/.local/bin on Windows) so the binary is on PATH.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Core Commands · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
