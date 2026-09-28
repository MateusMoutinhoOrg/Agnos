# `remove-available`

Deletes one available

```bash
agnos remove-available <Available> [--help] [--path <path>] [--quiet]
```

Removes adapters/availables/<name>/ whole. The standard available is refused: cmd/main/main.go imports it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Available` | string, required |  | the available to remove |

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
agnos remove-available lambda
```

Deps System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
