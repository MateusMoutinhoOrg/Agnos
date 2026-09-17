# `list-deps`

Lists the deps the embedded catalog can install

```bash
agnos list-deps [--path <path>] [--quiet]
```

One row per dep of the embedded catalog: the name, whether the project has the contract installed, the adapters filling it (or the catalog's default-adapter when none is), and what it provides.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos list-deps
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
