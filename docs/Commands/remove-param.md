# `remove-param`

Delete one declared query parameter from a route

```bash
agnos remove-param --route <route> [--path <path>] [--quiet] <name>
```

Drops one declared query parameter from a route, the exact inverse of add-param. The build renders only: dropping a parameter may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the parameter is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query parameter to drop |

```bash
agnos remove-param page --route list-users
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
