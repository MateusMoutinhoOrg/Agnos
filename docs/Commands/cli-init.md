# `cli-init`

Initializes the CLI layer for the project

```bash
agnos cli-init [--help] [--path <path>] [--quiet]
```

Installs the std and argv deps the CLI layer depends on, renders the "cli" asset group into the project, and calls build.

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
agnos cli-init
agnos cli-init --path ./my-project
```

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
