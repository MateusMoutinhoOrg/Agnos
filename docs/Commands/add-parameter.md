# `add-parameter`

Declare one value a route reads from the query string or the headers

```bash
agnos add-parameter --route <route> [--type <type>] [--font <font>...] [--required] [--default <default>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--example <example>...] [--position <position>] [--path <path>] [--quiet] <name>
```

Declares one value a route reads and runs build so the route's new.go and entries.go pick it up. It is read under its name from the first of its --font that brings it, converted to its --type and bound to Entries.<Name>; a missing --required one is answered 400 before the handler runs, and one with a --trigger only lets the route run when it matches.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the parameter |
| `--type` | string | `string` | the value type: string, integer, number, boolean, datetime, string-array or integer-array (defaults to string) |
| `--font` | string, repeatable |  | where the value is read from, in order: query, header or cookie (repeatable; defaults to query) |
| `--required` | boolean |  | answer 400 when no font brings the value (not for booleans or parameters with --default) |
| `--default` | string |  | the literal bound when no font brings the value (cannot be combined with --required) |
| `--trigger` | string |  | what the value has to match for the route to run at all, putting the parameter into what the route matches on |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix, text-prefix, suffix or regex — or starts-with, ends-with, exact, matches (defaults to equal; on a value, prefix and text-prefix are the same) |
| `--trigger-negate` | boolean |  | invert the trigger: the route runs when the value does not match it |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case |
| `--description` | string |  | help text shown for the parameter |
| `--example` | string, repeatable |  | an usage example for the parameter (repeatable) |
| `--position` | int | `-1` | zero-based index to insert the parameter at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query key or header name it is read under; its Entries field is the exported spelling of it |

```bash
agnos add-parameter page --route list-users --type number --default 1
agnos add-parameter authorization --route admin --font header --trigger Bearer --trigger-type prefix
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
