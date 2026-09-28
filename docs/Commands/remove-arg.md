# `remove-arg`

Remove a positional arg from a command's entries.yaml

```bash
agnos remove-arg <Name> --command <target> [--help] [--path <path>] [--quiet]
```

Drops one positional arg declaration from sandbox/internal/commands/<command>/entries.yaml and runs build so the command's new.go forgets it. Later args shift up.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the arg name |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the arg | — |
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
