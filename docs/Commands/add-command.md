# `add-command`

Scaffold a new command package in the project

```bash
agnos add-command <Name> --help <help> --category <category> [--path <path>] [--quiet]
```

Creates sandbox/internal/commands/<name>/ with a hand-written entries.yaml and a stub handler.go, then runs build so new.go and the dispatch pick it up. Refuses to overwrite an existing command.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the name of the new command (e.g. my-feature) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | string, required |  | one-line help text for the new command | — |
| `--category` | string, required |  | the category the new command is grouped under in help output | — |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-command my-feature
agnos add-command my-feature --path ./my-project
```

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
