# `remove-body-field`

Delete one property from a route's body json-schema

```bash
agnos remove-body-field <Name> --route <route> [--help] [--path <path>] [--quiet]
```

Drops one property of the body json-schema, named by the same dotted path add-body-field declared it with, and unlists it from its parent's required set. The build renders only: dropping a property may leave hand-written code referring to what is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the dotted path of the property to drop |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the property is declared on | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-body-field address.city --route create-user
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
