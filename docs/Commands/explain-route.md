# `explain-route`

Show which routes one request reaches, without a server

```bash
agnos explain-route <Method> <RequestPath> [--header <header>...] [--cookie <cookie>...] [--help] [--path <path>] [--quiet]
```

Walks the chain the way the generated dispatch does and prints, route by route, whether it runs for the request or why it is skipped — the segment count, a path that is missing, will not convert or fails its trigger, the method, a parameter trigger — and a 400 its binding would answer. What a handler does once it runs is its own code, so the last line says what the request ends on when none answers. Writes nothing.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Method` | string, required |  | the http method of the request |
| `RequestPath` | string, required |  | the request path, query string included |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--header` | string-array |  | a header the request carries, as key=value (repeatable) | — |
| `--cookie` | string-array |  | a cookie the request carries, as key=value (repeatable) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos explain-route GET /admin/users --header 'x-token=abc'
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
