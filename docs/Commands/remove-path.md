# `remove-path`

Delete one entry of a route's paths

```bash
agnos remove-path --route <route> [--path <path>] [--quiet] <id>
```

Drops one entry of the route's paths, the exact inverse of add-path, refusing the last one. The build renders only: dropping a path may leave hand-written code reading an Entries field that is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the path |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `id` | string, required |  | the id of the path to drop (its Entries field) |

```bash
agnos remove-path tenant --route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
