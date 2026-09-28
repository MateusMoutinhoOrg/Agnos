# `remove-command`

Delete a command package from the project

```bash
agnos remove-command <Name> [--help] [--path <path>] [--quiet]
```

Deletes sandbox/internal/commands/<name>/ (entries.yaml, new.go, handler.go and anything else inside) and runs build so the dispatch and help stop answering to it. The generated help command cannot be removed.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the command to delete (identifier or package name) |

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
agnos remove-command my-feature
agnos remove-command my-feature --path ./my-project
```

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
