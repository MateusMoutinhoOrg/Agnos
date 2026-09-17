# `server-init`

Add the http server layer to the project

```bash
agnos server-init [--path <path>] [--quiet]
```

Installs the deps the server layer needs, renders sandbox/internal/server, the routeio package and the built-in health route, and writes the start-server command. A project with no cli layer is given one first: a server needs a command that starts it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos server-init
agnos server-init --path ./my-project
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
