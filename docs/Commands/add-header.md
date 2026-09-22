# `add-header`

Declare a request header on a route

```bash
agnos add-header --route <route> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] [--identifier <identifier>] [--starts-with <starts-with>] <name>
```

Declares one request header on a route and runs build so the route's new.go picks it up. The name is the external spelling and is matched without regard to case; the dispatch answers 400 for a missing --required header or one outside --min/--max, before the handler runs.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the header |
| `--type` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description` | string |  | help text shown for the header |
| `--example` | string, repeatable |  | an usage example for the header (repeatable) |
| `--default` | string |  | the literal assigned when the header is absent (cannot be combined with --required) |
| `--required` | boolean |  | answer 400 when the header is not provided (not for booleans or headers with --default) |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the header at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--identifier` | string |  | match only when the request brings exactly this value, putting the field into what the route matches on |
| `--starts-with` | string |  | match only when the value the request brings begins with this, putting the field into what the route matches on |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the header name, matched without regard to case |

```bash
agnos add-header authorization --route create-user --required
agnos add-header x-retries --route create-user --type int --default 1 --max 5
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
