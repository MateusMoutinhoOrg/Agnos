# `remove-available`

Deletes one available

```bash
agnos remove-available [--path <path>] [--quiet] <available>
```

Removes adapters/availables/<name>/ whole. The standard available is refused: cmd/main/main.go imports it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `available` | string, required |  | the available to remove |

```bash
agnos remove-available lambda
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
