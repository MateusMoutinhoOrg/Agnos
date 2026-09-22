# `add-param`

Declare a query parameter on a route

```bash
agnos add-param --route <route> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] [--identifier <identifier>] [--starts-with <starts-with>] <name>
```

Declares one query-string parameter on a route and runs build so the route's new.go picks it up. --array collects every occurrence of the key, read back with GetStrings; the only other place it is accepted is the last segment of a route's paths, which takes the rest of the path.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the parameter |
| `--type` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description` | string |  | help text shown for the parameter |
| `--example` | string, repeatable |  | an usage example for the parameter (repeatable) |
| `--default` | string |  | the literal assigned when the parameter is absent (cannot be combined with --required) |
| `--required` | boolean |  | answer 400 when the parameter is not provided (not for booleans or parameters with --default) |
| `--array` | boolean |  | collect every occurrence into a []T field instead of a single value |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the parameter at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--identifier` | string |  | match only when the request brings exactly this value, putting the field into what the route matches on |
| `--starts-with` | string |  | match only when the value the request brings begins with this, putting the field into what the route matches on |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query key |

```bash
agnos add-param page --route list-users --type int --default 1 --min 1
agnos add-param tag --route list-users --array
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
