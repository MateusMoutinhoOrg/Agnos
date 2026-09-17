# `set-route`

Rewrite the route-level keys of a route.yaml

```bash
agnos set-route [--method <method>] [--help <help>] [--category <category>] [--long-description <long-description>] [--hidden] [--visible] [--path <path>] [--quiet] [--example <example>...] <route>
```

Overwrites method, help, category, long-description, hidden and examples on one route. Empty options leave the current value alone; --example appends.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--method`, `-m` | string |  | the http method the route answers |
| `--help` | string |  | one-line description of the route |
| `--category` | string |  | the heading the route is listed under in docs/Routes |
| `--long-description` | string |  | the paragraph docs/Routes prints under the route |
| `--hidden` | boolean |  | drop the route from docs/Routes, still dispatched |
| `--visible` | boolean |  | list the route again in docs/Routes |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--example` | string, repeatable |  | an usage example for the route (repeatable) |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route to edit (identifier or package name) |

```bash
agnos set-route create-user --method POST --example "curl -X POST localhost:8080/users"
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
