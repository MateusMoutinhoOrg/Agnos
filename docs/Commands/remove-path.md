# `remove-path`

Delete one entry of a route's paths

```bash
agnos remove-path <Id> --route <route> [--help] [--path <path>] [--quiet]
```

Drops one entry of the route's paths, the exact inverse of add-path, refusing the last one. The build renders only: dropping a path may leave hand-written code reading an Entries field that is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Id` | string, required |  | the id of the path to drop (its Entries field) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the path | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-path tenant --route create-user
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
