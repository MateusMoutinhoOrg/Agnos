# `rename-route`

Rename one route

```bash
agnos rename-route <Route> <Name> [--help] [--path <path>] [--quiet]
```

Moves sandbox/internal/routeslist/<route>/ to sandbox/internal/routeslist/<name>/, rewriting the package clause of every hand-written Go file, and runs build so new.go, entries.go and the server's route list follow it. A page is refused: remove-page and add-page own its html.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Route` | string, required |  | the route to rename (identifier or package name) |
| `Name` | string, required |  | the name it takes on |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
