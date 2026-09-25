# `add-route`

Declare a new http route

```bash
agnos add-route [--trigger <trigger>] [--trigger-type <trigger-type>] [--method <method>...] [--response-type <response-type>] --help <help> --category <category> [--path <path>] [--quiet] [--priority <priority>] <name>
```

Writes sandbox/internal/routeslist/<name>/route.yaml — one path, Route, reading the whole request path against the trigger — and a stub InternalPureHandler.go, then runs build so the route's new.go (the api.Route that lands in Server.Routes) and entries.go (the Entries its handler is handed) are generated. priority and response-type are always written.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--trigger` | string |  | what the whole request path is compared against; an equal, prefix or regex trigger is normalized to start with / (defaults to /<name>) |
| `--trigger-type` | string |  | how the trigger is compared against the whole request path: equal, prefix, suffix or regex (defaults to equal) |
| `--method`, `-m` | string, repeatable |  | an http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD or OPTIONS (repeatable; defaults to GET) |
| `--response-type` | string |  | the Content-Type every response of the route carries (defaults to application/json) |
| `--help` | string, required |  | one-line description of the route |
| `--category` | string, required |  | the heading the route is listed under in docs/Routes |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--priority` | int | `0` | the rung this route runs on when several match one request: lowest first, and a route that writes no status hands the request on |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route name (becomes the directory sandbox/internal/routeslist/<name> and its Go package) |

```bash
agnos add-route create-user --trigger /users --method POST --help "Create a user" --category Users
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
