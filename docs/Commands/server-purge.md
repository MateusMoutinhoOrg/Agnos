# `server-purge`

Remove the http server layer and every route in it

```bash
agnos server-purge [--help] [--path <path>] [--quiet]
```

Drops sandbox/internal/{server,routes} and the start-server command, then re-renders. The cli layer and the installed deps — the OpinionatedAgnosServer lib among them — are left in place.

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
agnos server-purge
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
