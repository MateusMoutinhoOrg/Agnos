# `add-route`

Declare a new http route

```bash
agnos add-route <Name> [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--pattern <pattern>] [--method <method>...] [--middleware] [--priority <priority>] [--before <before>] [--after <after>] [--response-type <response-type>] [--summary <summary>] [--category <category>] [--dir <dir>] [--help] [--path <path>] [--quiet]
```

Writes <name>/route.yaml and a stub handler.go under sandbox/internal/routes — or under the folder --dir names there, e.g. routes/admin/<name> — then runs build so the route's new.go (the api.Route that lands in Server.Routes) and input.go (the Input its handler is handed) are generated. A directory is a route by holding a route.yaml, at any depth; a name is unique across every folder. The paths come from --pattern, or from one path reading the whole request path against --trigger, named after its words (/api/products reads into ApiProducts). A --middleware declines by default and sits in front of the routes on the default rung; --before and --after place a route next to another one. priority and response-type are always written.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route name (becomes the directory sandbox/internal/routes/<name> and its Go package) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--trigger` | string |  | what the whole request path is compared against; an equal, prefix or text-prefix trigger is normalized to start with / (defaults to /<name>, or / for a --middleware) | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (the value or anything under it, segment by segment), text-prefix, suffix or regex — or starts-with, ends-with, exact, matches (defaults to equal, or prefix for a --middleware) | — |
| `--trigger-negate` | boolean |  | invert the trigger: the route runs for every path that does not match it | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--pattern` | string |  | the url shape the route answers, compiled into its paths: literal segments, {name}, {name:integer\|number\|uuid} and a last {*rest}; without {*rest} the segment count is fixed (excludes --trigger and --trigger-type) | — |
| `--method`, `-m` | string-array |  | an http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, or ANY for every one (repeatable; defaults to GET, or ANY for a --middleware) | — |
| `--middleware` | boolean |  | declare a middleware: ANY method, a prefix trigger on / unless --trigger says otherwise, priority 10, and a stub handler that hands the request on | — |
| `--priority` | integer | `-1` | the rung this route runs on when several match one request, lowest first; a route that answers nothing hands the request on (defaults to 100, or 10 for a --middleware) | — |
| `--before` | string |  | land one rung below the route named, so it runs first (excludes --priority) | — |
| `--after` | string |  | land one rung above the route named, so it runs next (excludes --priority) | — |
| `--response-type` | string |  | the Content-Type every response of the route carries (defaults to application/json, or text/plain for a --middleware) | — |
| `--summary` | string |  | one-line description of the route | — |
| `--category` | string |  | the heading the route is listed under in docs/Routes (defaults to Routes, or Middleware for a --middleware) | — |
| `--dir` | string |  | the folder under sandbox/internal/routes the route lands in, e.g. admin puts it in routes/admin/<name> (defaults to the top); a directory holding a route.yaml is a route, whatever folder holds it | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-route create-user --trigger /users --method POST --summary "Create a user" --category Users
agnos add-route get-article --pattern '/articles/{article:integer}'
agnos add-route admin-guard --middleware --trigger /admin
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
