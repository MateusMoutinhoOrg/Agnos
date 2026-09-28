# `server-purge`

Remove the http server layer and every route in it

```bash
agnos server-purge [--help] [--path <path>] [--quiet]
```

Drops sandbox/internal/{server,routeslist,routeio} and the start-server command, then re-renders. The cli layer and the installed deps are left in place.

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
agnos server-purge
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
