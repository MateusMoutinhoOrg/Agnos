# `cli-purge`

Removes the CLI layer from the project

```bash
agnos cli-purge [--path <path>] [--quiet]
```

Removes every file the "cli" asset group installs and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos cli-purge
agnos cli-purge --path ./my-project
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
