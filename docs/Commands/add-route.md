# `add-route`

Declare a new http route

```bash
agnos add-route [--trigger <trigger>] [--method <method>] --help <help> --category <category> [--path <path>] [--quiet] [--priority <priority>] [--starts-with] <name>
```

Writes sandbox/internal/routes/<name>/route.yaml and a stub handler.go, then runs build so the route's new.go — the api.Route that lands in Server.Routes — is generated. The trigger is normalized to start with /, and defaults to /<name>.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--trigger` | string |  | the first literal segment of the path, always starting with / (defaults to /<name>) |
| `--method`, `-m` | string | `GET` | the http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD or OPTIONS |
| `--help` | string, required |  | one-line description of the route |
| `--category` | string, required |  | the heading the route is listed under in docs/Routes |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--priority` | int | `0` | the rung this route runs on when several match one request: lowest first, and a route that writes no status hands the request on |
| `--starts-with` | boolean |  | make the trigger a prefix, so the route answers every path beginning with it rather than that one path |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route name (becomes the directory sandbox/internal/routes/<name> and its Go package) |

```bash
agnos add-route create-user --trigger /users --method POST --help "Create a user" --category Users
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
