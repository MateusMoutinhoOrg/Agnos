# `remove-body-field`

Delete one property from a route's body json-schema

```bash
agnos remove-body-field --route <route> [--path <path>] [--quiet] <name>
```

Drops one property of the body json-schema, named by the same dotted path add-body-field declared it with, and unlists it from its parent's required set. The build renders only: dropping a property may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the property is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the dotted path of the property to drop |

```bash
agnos remove-body-field address.city --route create-user
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
