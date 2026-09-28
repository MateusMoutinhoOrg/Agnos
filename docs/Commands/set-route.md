# `set-route`

Rewrite the route-level keys of a route.yaml

```bash
agnos set-route <Route> [--method <method>...] [--response-type <response-type>] [--help <help>] [--category <category>] [--long-description <long-description>] [--hidden] [--visible] [--priority <priority>] [--before <before>] [--after <after>] [--segments <segments>] [--clear <clear>...] [--example <example>...] [--path <path>] [--quiet]
```

Overwrites methods, response-type, priority, segments, help, category, long-description, hidden and examples on one route. Empty options leave the current value alone; --method replaces the whole list; --example appends; --before and --after place the route one rung from another; --clear takes a key off.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Route` | string, required |  | the route to edit (identifier or package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--method`, `-m` | string-array |  | an http method the route answers, or ANY for every one (repeatable; replaces the whole list) | — |
| `--response-type` | string |  | the Content-Type every response of the route carries | — |
| `--help` | string |  | one-line description of the route | — |
| `--category` | string |  | the heading the route is listed under in docs/Routes | — |
| `--long-description` | string |  | the paragraph docs/Routes prints under the route | — |
| `--hidden` | boolean |  | drop the route from docs/Routes, still dispatched | — |
| `--visible` | boolean |  | list the route again in docs/Routes | — |
| `--priority` | integer | `-1` | the rung this route runs on when several match one request: lowest first, and a route that answers nothing hands the request on | — |
| `--before` | string |  | move to one rung below the route named, so it runs first (excludes --priority) | — |
| `--after` | string |  | move to one rung above the route named, so it runs next (excludes --priority) | — |
| `--segments` | integer | `-1` | how many segments the request path has to have for the route to run | — |
| `--clear` | string-array |  | a key to take off: segments (repeatable) | — |
| `--example` | string-array |  | an usage example for the route (repeatable) | — |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos set-route create-user --method POST --example "curl -X POST localhost:8080/users"
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
