# `local-install`

Builds the project and installs it locally

```bash
agnos local-install [--path <path>] [--quiet]
```

Runs build over the project, then compiles ./cmd/main into /usr/local/bin/<project-name> (~/.local/bin on Windows) so the binary is on PATH.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

Core Commands · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
