# `add-available`

Declares one further available

```bash
agnos add-available [--path <path>] [--quiet] <available>
```

Creates adapters/availables/<name>/available.yaml as a copy of the standard selection, so it starts filling every field, and build generates its new.go. Point it at another adapter with set-adapter --available.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `available` | string, required |  | the name of the new available (it becomes one directory under adapters/availables/) |

```bash
agnos add-available lambda
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
