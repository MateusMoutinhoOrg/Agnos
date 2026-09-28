# `add-available`

Declares one further available

```bash
agnos add-available <Available> [--help] [--path <path>] [--quiet]
```

Creates adapters/availables/<name>/available.yaml as a copy of the standard selection, so it starts filling every field, and build generates its new.go. Point it at another adapter with set-adapter --available.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Available` | string, required |  | the name of the new available (it becomes one directory under adapters/availables/) |

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
agnos add-available lambda
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
