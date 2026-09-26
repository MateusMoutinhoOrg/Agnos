# `rename-route`

Rename one route

```bash
agnos rename-route [--path <path>] [--quiet] <route> <name>
```

Moves sandbox/internal/routeslist/<route>/ to sandbox/internal/routeslist/<name>/, rewriting the package clause of every hand-written Go file, and runs build so new.go, entries.go and the server's route list follow it. A page is refused: remove-page and add-page own its html.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route to rename (identifier or package name) |
| `name` | string, required |  | the name it takes on |

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
