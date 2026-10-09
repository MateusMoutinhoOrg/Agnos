# `rename-route`

Rename one route

```bash
agnos rename-route <Route> <Name> [--dir <dir>] [--help] [--path <path>] [--quiet]
```

Moves the route's directory to <name>/ in the folder it sits in — or in the one --dir names, / for the top — rewriting the package clause of every hand-written Go file, drops every folder that leaves empty, and runs build so generated.new.go, generated.input.go and the server's route list follow it. <name> may be the current one when only the folder changes. A page is refused: remove-page and add-page own its html.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Route` | string, required |  | the route to rename (identifier or package name) |
| `Name` | string, required |  | the name it takes on |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--dir` | string |  | the folder under sandbox/internal/routes the route moves to, / for the top (defaults to the folder it sits in) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
