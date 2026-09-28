# `show-command`

Print one command's declaration as a tree: args, flags and the middlewares in front of it

```bash
agnos show-command <Target> [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Target` | string, required |  | the command to show (its name or a verb it answers to) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
