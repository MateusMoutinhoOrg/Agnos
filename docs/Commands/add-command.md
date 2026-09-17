# `add-command`

Scaffold a new command package in the project

```bash
agnos add-command --help <help> --category <category> [--path <path>] [--quiet] <name>
```

Creates sandbox/internal/commands/<name>/ with a hand-written entries.yaml and a stub handler.go, then runs build so new.go and the dispatch pick it up. Refuses to overwrite an existing command.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--help` | string, required |  | one-line help text for the new command |
| `--category` | string, required |  | the category the new command is grouped under in help output |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new command (e.g. my-feature) |

```bash
agnos add-command my-feature
agnos add-command my-feature --path ./my-project
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
