# `remove-parameter`

Delete one entry of a route's parameters

```bash
agnos remove-parameter --route <route> [--path <path>] [--quiet] <name>
```

Drops one entry of the route's parameters, the exact inverse of add-parameter. The build renders only: dropping a parameter may leave hand-written code reading an Entries field that is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the parameter |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the key of the parameter to drop |

```bash
agnos remove-parameter page --route list-users
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
