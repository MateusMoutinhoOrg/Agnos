# `set-path`

Rewrite one entry of a route's paths

```bash
agnos set-path <Id> --route <route> [--rename <rename>] [--start <start>] [--end <end>] [--type <type>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

Rewrites one entry of the route's paths in place and runs build. It is add-path applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Id` | string, required |  | the id of the path to rewrite |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the path | — |
| `--rename` | string |  | the id the path takes on | — |
| `--start` | string |  | the new index of the first segment | — |
| `--end` | string |  | the new index of the last segment, -1 for the last one | — |
| `--type` | string |  | the new type: string, integer, number or uuid | — |
| `--trigger` | string |  | the new value the slice has to read as | — |
| `--trigger-type` | string |  | the new comparison: equal, prefix, text-prefix, suffix or regex — or starts-with, ends-with, exact, matches | — |
| `--trigger-negate` | boolean |  | invert the trigger (take it off with --clear trigger-negate) | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case (take it off with --clear trigger-ignore-case) | — |
| `--description` | string |  | the new help text | — |
| `--clear` | string-array |  | a key to take off: trigger, trigger-negate, trigger-ignore-case, type or description (repeatable) | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos set-path tenant --route create-user --end 2
agnos set-path Route --route logger --trigger-type prefix
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
