# `remove-flag`

Remove a flag from a command's entries.yaml

```bash
agnos remove-flag <Name> --command <target> [--help] [--path <path>] [--quiet]
```

Drops one flag declaration (matched by its name or by one of its identifiers) from sandbox/internal/commands/<command>/entries.yaml and runs build so the command's new.go forgets it.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the flag name (or one of its identifiers, e.g. --out) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the flag | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-flag output --command exec
agnos remove-flag --out --command exec
```

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
