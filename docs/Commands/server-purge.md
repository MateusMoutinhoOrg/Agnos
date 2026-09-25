# `server-purge`

Remove the http server layer and every route in it

```bash
agnos server-purge [--path <path>] [--quiet]
```

Drops sandbox/internal/{server,routeslist,routeio} and the start-server command, then re-renders. The cli layer and the installed deps are left in place.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos server-purge
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
