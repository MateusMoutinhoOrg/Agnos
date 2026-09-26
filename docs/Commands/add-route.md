# `add-route`

Declare a new http route

```bash
agnos add-route [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--pattern <pattern>] [--method <method>...] [--middleware] [--priority <priority>] [--before <before>] [--after <after>] [--phase <phase>] [--response-type <response-type>] [--help <help>] [--category <category>] [--path <path>] [--quiet] <name>
```

Writes sandbox/internal/routeslist/<name>/route.yaml and a stub InternalPureHandler.go, then runs build so the route's new.go (the api.Route that lands in Server.Routes) and entries.go (the Entries its handler is handed) are generated. The paths come from --pattern, or from one path, Route, reading the whole request path against --trigger. A --middleware declines by default and sits in front of the routes on the default rung; --before and --after place a route next to another one. priority and response-type are always written.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--trigger` | string |  | what the whole request path is compared against; an equal, prefix or text-prefix trigger is normalized to start with / (defaults to /<name>, or / for a --middleware) |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (the value or anything under it, segment by segment), text-prefix, suffix or regex — or starts-with, ends-with, exact, matches (defaults to equal, or prefix for a --middleware) |
| `--trigger-negate` | boolean |  | invert the trigger: the route runs for every path that does not match it |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case |
| `--pattern` | string |  | the url shape the route answers, compiled into its paths: literal segments, {name}, {name:integer\|number\|uuid} and a last {*rest}; without {*rest} the segment count is fixed (excludes --trigger and --trigger-type) |
| `--method`, `-m` | string, repeatable |  | an http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, or ANY for every one (repeatable; defaults to GET, or ANY for a --middleware) |
| `--middleware` | boolean |  | declare a middleware: ANY method, a prefix trigger on / unless --trigger says otherwise, priority 10, and a stub handler that hands the request on |
| `--priority` | int |  | the rung this route runs on when several match one request, lowest first; a route that answers nothing hands the request on (defaults to 100, or 10 for a --middleware) |
| `--before` | string |  | land one rung below the route named, so it runs first (excludes --priority) |
| `--after` | string |  | land one rung above the route named, so it runs next (excludes --priority) |
| `--phase` | string |  | when the route runs: before (the chain, the default) or after, once the request has been answered and without being able to change the answer |
| `--response-type` | string |  | the Content-Type every response of the route carries (defaults to application/json, or text/plain for a --middleware) |
| `--help` | string |  | one-line description of the route |
| `--category` | string |  | the heading the route is listed under in docs/Routes (defaults to Routes, or Middleware for a --middleware) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route name (becomes the directory sandbox/internal/routeslist/<name> and its Go package) |

```bash
agnos add-route create-user --trigger /users --method POST --help "Create a user" --category Users
agnos add-route get-article --pattern '/articles/{article:integer}'
agnos add-route admin-guard --middleware --trigger /admin
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
