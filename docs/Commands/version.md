# `version`

Print the installed version

```bash
agnos version [--help] [--path <path>] [--quiet]
```

Prints the current version of the installed binary and exits.

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
agnos version
```

Info · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
