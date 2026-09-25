# `remove-route`

Delete one declared route

```bash
agnos remove-route [--path <path>] [--quiet] <name>
```

Removes sandbox/internal/routeslist/<name>/ whole and re-renders the dispatch. The build renders only: dropping a route may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route to delete (identifier or package name) |

```bash
agnos remove-route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
