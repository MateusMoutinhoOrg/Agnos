# `remove-header`

Delete one declared header from a route

```bash
agnos remove-header --route <route> [--path <path>] [--quiet] <name>
```

Drops one declared header from a route, the exact inverse of add-header. The build renders only: dropping a header may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the header is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the header to drop |

```bash
agnos remove-header authorization --route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
