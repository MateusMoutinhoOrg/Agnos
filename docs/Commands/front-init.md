# `front-init`

Add the html front layer to the project

```bash
agnos front-init [--path <path>] [--quiet]
```

Installs the deps the front layer needs, renders sandbox/internal/frontio and writes, once, the frontend route serving every file of assets/frontend and that tree's index.html and 404.html. A project with no server layer is given one first: the front is answered over http.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos front-init
agnos front-init --path ./my-project
```

Front System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
