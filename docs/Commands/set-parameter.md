# `set-parameter`

Rewrite one entry of a route's parameters

```bash
agnos set-parameter <Name> --route <route> [--rename <rename>] [--type <type>] [--source <source>...] [--required] [--default <default>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--example <example>...] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

Rewrites one entry of the route's parameters in place and runs build. It is add-parameter applied to a declaration that already exists: the keys given are written over the ones there, --clear takes one off, and the result goes through the same constructor.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the key of the parameter to rewrite |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that declares the parameter | — |
| `--rename` | string |  | the key the parameter takes on | — |
| `--type` | string |  | the new value type: string, integer, number, boolean, datetime, string-array or integer-array | — |
| `--source` | string-array |  | where the value is read from, in order: query, header or cookie (repeatable; replaces the whole list) | — |
| `--required` | boolean |  | answer 400 when no source brings the value | — |
| `--default` | string |  | the new literal bound when no source brings the value | — |
| `--trigger` | string |  | the new value the parameter has to match for the route to run | — |
| `--trigger-type` | string |  | the new comparison: equal, prefix, text-prefix, suffix or regex — or starts-with, ends-with, exact, matches | — |
| `--trigger-negate` | boolean |  | invert the trigger (take it off with --clear trigger-negate) | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case (take it off with --clear trigger-ignore-case) | — |
| `--description` | string |  | the new help text | — |
| `--example` | string-array |  | an extra usage example (repeatable) | — |
| `--clear` | string-array |  | a key to take off: description, examples, default, required, trigger, trigger-negate or trigger-ignore-case (repeatable) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos set-parameter page --route list-users --type number --default 1
agnos set-parameter token --route admin --source header --clear trigger
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
