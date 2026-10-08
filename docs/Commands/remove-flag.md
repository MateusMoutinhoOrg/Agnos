# `remove-flag`

Remove a flag from a command's command.yaml

```bash
agnos remove-flag <Name> --command <command> [--help] [--path <path>] [--quiet]
```

Drops one flag declaration (matched by its name, its id or one of its keys) from sandbox/internal/commands/<command>/command.yaml and runs build so the command's new.go and input.go follow it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the flag name, its id, or one of its keys, e.g. --out |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that owns the flag | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-flag output --command exec
agnos remove-flag out --command exec
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
