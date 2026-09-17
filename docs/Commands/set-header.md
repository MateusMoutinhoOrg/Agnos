# `set-header`

Rewrite one declared request header of a route

```bash
agnos set-header --route <route> [--rename <rename>] [--type <type>] [--description <description>] [--default <default>] [--required] [--min <min>] [--max <max>] [--example <example>...] [--clear <clear>...] [--path <path>] [--quiet] <name>
```

Rewrites one declared header in place and runs build. It is add-header applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor — so adding a bound that was forgotten never means removing the header and declaring it again.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the header is declared on |
| `--rename` | string |  | the header name it answers to from now on |
| `--type` | string |  | the value type: string, boolean, int or float |
| `--description` | string |  | help text shown for the header |
| `--default` | string |  | the literal bound when the header is absent |
| `--required` | boolean |  | refuse a request that does not carry the header |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--example` | string, repeatable |  | a value the header takes, for the docs (repeatable; appended) |
| `--clear` | string, repeatable |  | a key to take off again: description, examples, default, required, min or max (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the header to edit, by the name it is declared under |

```bash
agnos set-header authorization --route create-user --required --description 'the bearer token'
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
