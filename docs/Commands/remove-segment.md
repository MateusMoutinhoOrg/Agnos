# `remove-segment`

Delete one segment from a route's path

```bash
agnos remove-segment --route <route> [--path <path>] [--quiet] <name>
```

Drops one segment from the route's paths, the exact inverse of add-segment: a capture by its name, a literal by the identifier it spells. The build renders only: dropping a segment may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the segment is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the segment to drop: the captured segment's name, or the identifier of a literal one |

```bash
agnos remove-segment tenant --route create-user
agnos remove-segment /users --route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
