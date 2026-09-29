# `remove-parameter`

Delete one entry of a route's parameters

```bash
agnos remove-parameter <Name> --route <route> [--help] [--path <path>] [--quiet]
```

Drops one entry of the route's parameters, the exact inverse of add-parameter. The build renders only: dropping a parameter may leave hand-written code reading an Entries field that is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the key of the parameter to drop |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the parameter | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-parameter page --route list-users
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
