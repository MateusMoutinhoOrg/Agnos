# `explain-route`

Show which routes one request reaches, without a server

```bash
agnos explain-route [--header <header>...] [--cookie <cookie>...] [--path <path>] [--quiet] <method> <request-path>
```

Walks the chain the way the generated dispatch does and prints, route by route, whether it runs for the request or why it is skipped — the segment count, a path that is missing, will not convert or fails its trigger, the method, a parameter trigger — and a 400 its binding would answer. What a handler does once it runs is its own code, so the last line says what the request ends on when none answers. Writes nothing.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--header` | string, repeatable |  | a header the request carries, as key=value (repeatable) |
| `--cookie` | string, repeatable |  | a cookie the request carries, as key=value (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `method` | string, required |  | the http method of the request |
| `request-path` | string, required |  | the request path, query string included |

```bash
agnos explain-route GET /admin/users --header 'x-token=abc'
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
