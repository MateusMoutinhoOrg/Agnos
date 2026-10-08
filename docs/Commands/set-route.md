# `set-route`

Rewrite the route-level keys of a route.yaml

```bash
agnos set-route <Name> [--method <method>...] [--response-type <response-type>] [--summary <summary>] [--category <category>] [--description <description>] [--hidden] [--visible] [--priority <priority>] [--before <before>] [--after <after>] [--segments <segments>] [--clear <clear>...] [--example <example>...] [--help] [--path <path>] [--quiet]
```

Overwrites methods, response-type, priority, segments, summary, category, description, hidden and examples on one route. Empty options leave the current value alone; --method replaces the whole list; --example appends; --before and --after place the route one rung from another; --clear takes a key off.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route to edit (identifier or package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--method`, `-m` | string-array |  | an http method the route answers, or ANY for every one (repeatable; replaces the whole list) | — |
| `--response-type` | string |  | the Content-Type every response of the route carries | — |
| `--summary` | string |  | one-line description of the route | — |
| `--category` | string |  | the heading the route is listed under in docs/Routes | — |
| `--description` | string |  | the paragraph docs/Routes prints under the route | — |
| `--hidden` | boolean |  | drop the route from docs/Routes, still dispatched | — |
| `--visible` | boolean |  | list the route again in docs/Routes | — |
| `--priority` | integer | `-1` | the rung this route runs on when several match one request: lowest first, and a route that answers nothing hands the request on | — |
| `--before` | string |  | move to one rung below the route named, so it runs first (excludes --priority) | — |
| `--after` | string |  | move to one rung above the route named, so it runs next (excludes --priority) | — |
| `--segments` | integer | `-1` | how many segments the request path has to have for the route to run | — |
| `--clear` | string-array |  | a key to take off: segments (repeatable) | — |
| `--example` | string-array |  | an usage example for the route (repeatable) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos set-route create-user --method POST --example "curl -X POST localhost:8080/users"
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
