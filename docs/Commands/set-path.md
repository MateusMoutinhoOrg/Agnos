# `set-path`

Rewrite one entry of a route's paths

```bash
agnos set-path --route <route> [--rename <rename>] [--start <start>] [--end <end>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--description <description>] [--clear <clear>...] [--path <path>] [--quiet] <id>
```

Rewrites one entry of the route's paths in place and runs build. It is add-path applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the path |
| `--rename` | string |  | the id the path takes on |
| `--start` | string |  | the new index of the first segment |
| `--end` | string |  | the new index of the last segment, -1 for the last one |
| `--trigger` | string |  | the new value the slice has to read as |
| `--trigger-type` | string |  | the new comparison: equal, prefix, suffix or regex |
| `--description` | string |  | the new help text |
| `--clear` | string, repeatable |  | a key to take off: trigger or description (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `id` | string, required |  | the id of the path to rewrite |

```bash
agnos set-path tenant --route create-user --end 2
agnos set-path Route --route logger --trigger-type prefix
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
