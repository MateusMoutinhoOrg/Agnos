# `remove-path`

Delete one entry of a route's paths

```bash
agnos remove-path <Name> --route <route> [--help] [--path <path>] [--quiet]
```

Drops one entry of the route's paths, the exact inverse of add-path, refusing the last one. The build renders only: dropping a path may leave hand-written code reading an Input field that is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the id of the path to drop (its Input field) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the path | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-path tenant --route create-user
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
