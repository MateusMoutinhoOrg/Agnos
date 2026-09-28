# `remove-route`

Delete one declared route

```bash
agnos remove-route <Name> [--help] [--path <path>] [--quiet]
```

Removes sandbox/internal/routeslist/<name>/ whole and re-renders the dispatch. The build renders only: dropping a route may leave hand-written code referring to what is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route to delete (identifier or package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos remove-route create-user
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
