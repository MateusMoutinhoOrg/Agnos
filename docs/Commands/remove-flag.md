# `remove-flag`

Remove a flag from a command's entries.yaml

```bash
agnos remove-flag --command <command> [--path <path>] [--quiet] <name>
```

Drops one flag declaration (matched by its name or by one of its identifiers) from sandbox/internal/commands/<command>/entries.yaml and runs build so the command's new.go forgets it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the flag |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the flag name (or one of its identifiers, e.g. --out) |

```bash
agnos remove-flag output --command exec
agnos remove-flag --out --command exec
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
