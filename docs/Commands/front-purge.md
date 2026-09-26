# `front-purge`

Remove the html front layer from the project

```bash
agnos front-purge [--path <path>] [--quiet]
```

Removes sandbox/internal/generated/frontio and the frontend route, then rebuilds. assets/frontend/ is left untouched: every file there is the project's content, so front-init puts the route back over it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos front-purge
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
