# `add-segment`

Add a segment to a route's path

```bash
agnos add-segment --route <route> [--identifier <identifier>] [--type <type>] [--description <description>] [--example <example>...] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] [--starts-with <starts-with>] [<name>]
```

Appends one segment to the route's paths and runs build so the route's new.go picks it up. With --identifier the segment is a literal, normalized to start with /; with a name it is a capture, which is always required and is bound under that name already converted. --array makes that capture take every segment left in the path, which only the last segment of a route may do.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the segment |
| `--identifier` | string |  | declare a literal segment instead of a capture, always starting with / |
| `--type` | string | `string` | the value type of a captured segment: string, boolean, int or float |
| `--description` | string |  | help text shown for the captured segment |
| `--example` | string, repeatable |  | an usage example for the segment (repeatable) |
| `--array` | boolean |  | take every segment left in the path into a []T field (the last segment only) |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the segment at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--starts-with` | string |  | spell a literal segment the path only has to begin with, so the route answers every path under it |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string |  | the captured segment's name (omitted when --identifier declares a literal segment) |

```bash
agnos add-segment --route create-user --identifier /users
agnos add-segment tenant --route create-user --description "the tenant the user belongs to"
agnos add-segment page --route list-users --type int --min 1
agnos add-segment rest --route static --array
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
