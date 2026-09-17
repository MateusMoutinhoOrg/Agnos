# `front-purge`

Remove the html front layer from the project

```bash
agnos front-purge [--path <path>] [--quiet]
```

Removes sandbox/internal/pageio, the static route and the route of every declared page, then rebuilds. assets/frontend/ is left untouched: pages, styles and scripts are hand-written content, so front-init followed by add-page puts the routes back over it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos front-purge
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
