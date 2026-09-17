# `set-segment`

Rewrite one declared segment of a route's path

```bash
agnos set-segment --route <route> [--rename <rename>] [--identifier <identifier>] [--type <type>] [--description <description>] [--array] [--min <min>] [--max <max>] [--example <example>...] [--clear <clear>...] [--path <path>] [--quiet] <name>
```

Rewrites one segment of the route's path in place and runs build. It is add-segment applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor — so editing a segment and declaring one leave the same bytes. --identifier makes the segment a literal, whatever it was before; every other key edits a capture.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the segment is declared on |
| `--rename` | string |  | the name the captured segment answers to from now on |
| `--identifier` | string |  | make the segment a literal spelling this, whatever it was before |
| `--type` | string |  | the value type of the capture: string, boolean, int or float |
| `--description` | string |  | help text shown for the capture |
| `--array` | boolean |  | make the capture take every segment left in the path (last segment only) |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--example` | string, repeatable |  | a value the capture takes, for the docs (repeatable; appended) |
| `--clear` | string, repeatable |  | a key to take off again: description, examples, required, array, min or max (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the segment to edit: a capture by its name, a literal by the identifier it spells |

```bash
agnos set-segment tenant --route create-user --type int --description 'the tenant id'
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
