# `remove-arg`

Remove an arg from a command's command.yaml

```bash
agnos remove-arg <Name> --command <command> [--help] [--path <path>] [--quiet]
```

Drops one arg declaration from sandbox/internal/commands/<command>/command.yaml and runs build so the command's new.go and entries.go follow it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the arg name |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that owns the arg | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-arg file --command exec
```

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
