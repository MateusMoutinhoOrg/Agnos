# `rename-command`

Rename a command: its package, and the verb it answered to by its name

```bash
agnos rename-command <Target> <Name> [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Target` | string, required |  | the command to rename (its name or a verb it answers to) |
| `Name` | string, required |  | the name it takes on |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
