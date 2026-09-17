# `remove-arg`

Remove a positional arg from a command's entries.yaml

```bash
agnos remove-arg --command <command> [--path <path>] [--quiet] <name>
```

Drops one positional arg declaration from sandbox/internal/commands/<command>/entries.yaml and runs build so the command's new.go forgets it. Later args shift up.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the arg |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the arg name |

```bash
agnos remove-arg file --command exec
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
