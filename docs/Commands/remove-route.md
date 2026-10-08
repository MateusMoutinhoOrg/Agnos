# `remove-route`

Delete one declared route

```bash
agnos remove-route <Name> [--help] [--path <path>] [--quiet]
```

Removes the route's directory whole, in whatever folder of sandbox/internal/routes it sits, and every folder that leaves empty, then re-renders the dispatch. A directory holding another route is refused. The build renders only: dropping a route may leave hand-written code referring to what is gone.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route to delete (identifier or package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-route create-user
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
