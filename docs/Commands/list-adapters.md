# `list-adapters`

Lists the adapters of the catalog and of the project

```bash
agnos list-adapters [--path <path>] [--quiet]
```

One row per adapter: the name, the dep it fills, whether it is installed, the availables binding it, and what backs it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos list-adapters
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
