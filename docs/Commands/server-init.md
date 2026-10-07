# `server-init`

Add the http server layer to the project

```bash
agnos server-init [--help] [--path <path>] [--quiet]
```

Installs the deps the server layer needs — the OpinatedAgnosServer lib, its request chain, among them — renders sandbox/internal/server and the built-in health route, and writes the start-server command. A project with no cli layer is given one first: a server needs a command that starts it.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos server-init
agnos server-init --path ./my-project
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
