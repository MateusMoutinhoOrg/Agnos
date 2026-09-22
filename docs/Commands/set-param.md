# `set-param`

Rewrite one declared query parameter of a route

```bash
agnos set-param --route <route> [--rename <rename>] [--type <type>] [--description <description>] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--example <example>...] [--clear <clear>...] [--path <path>] [--quiet] [--identifier <identifier>] [--starts-with <starts-with>] <name>
```

Rewrites one declared query parameter in place and runs build. It is add-param applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor — so adding a bound that was forgotten never means removing the parameter and declaring it again.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the parameter is declared on |
| `--rename` | string |  | the query key it answers to from now on |
| `--type` | string |  | the value type: string, boolean, int or float |
| `--description` | string |  | help text shown for the parameter |
| `--default` | string |  | the literal bound when the parameter is absent |
| `--required` | boolean |  | refuse a request that does not carry the parameter |
| `--array` | boolean |  | collect every occurrence of the key into a []T field |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--example` | string, repeatable |  | a value the parameter takes, for the docs (repeatable; appended) |
| `--clear` | string, repeatable |  | a key to take off again: description, examples, default, required, array, min or max (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--identifier` | string |  | match only when the request brings exactly this value, putting the field into what the route matches on |
| `--starts-with` | string |  | match only when the value the request brings begins with this, putting the field into what the route matches on |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query parameter to edit, by the key it is declared under |

```bash
agnos set-param page --route list-users --type int --min 1 --max 50
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
