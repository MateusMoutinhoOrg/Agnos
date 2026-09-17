# `cli-init`

Initializes the CLI layer for the project

```bash
agnos cli-init [--path <path>] [--quiet]
```

Installs the std and argv deps the CLI layer depends on, renders the "cli" asset group into the project, and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos cli-init
agnos cli-init --path ./my-project
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
