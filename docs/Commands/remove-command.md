# `remove-command`

Delete a command package from the project

```bash
agnos remove-command [--path <path>] [--quiet] <name>
```

Deletes sandbox/internal/commands/<name>/ (entries.yaml, new.go, handler.go and anything else inside) and runs build so the dispatch and help stop answering to it. The generated help command cannot be removed.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the command to delete (identifier or package name) |

```bash
agnos remove-command my-feature
agnos remove-command my-feature --path ./my-project
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
