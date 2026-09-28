# `explain-command`

Run a command line against the declared commands without running any

```bash
agnos explain-command [Argv…] [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Argv` | string, repeatable |  | the command line to explain, after a bare -- so its flags are not read as this command's |

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
