# `add-path`

Add one slice of the request path to a route

```bash
agnos add-path <Name> --route <route> [--start <start>] [--end <end>] [--type <type>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--position <position>] [--help] [--path <path>] [--quiet]
```

Inserts one entry into the route's paths and runs build so the route's generated.new.go and generated.input.go pick it up. A path reads the request segments from --start to --end (both inclusive, -1 the last one) as '/' followed by them joined by '/', binds that text to Input.<Id>, and — with --trigger — only lets the route run when the text matches.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the id of the path, the Input field its slice binds to (normalized to an exported Go name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the path | — |
| `--start` | string |  | the index of the first request segment the slice reads (defaults to 0) | — |
| `--end` | string |  | the index of the last request segment the slice reads, -1 for the last one (defaults to -1) | — |
| `--type` | string |  | what the slice converts to: string (the default), integer, number or uuid; anything but string reads one segment (--start equal to --end), and a segment that does not convert makes the route a non-match | — |
| `--trigger` | string |  | what the slice has to read as for the route to run; without it the path is a plain capture | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (segment by segment), text-prefix, suffix or regex — or starts-with, ends-with, exact, matches (defaults to equal) | — |
| `--trigger-negate` | boolean |  | invert the trigger: the route runs when the slice does not match it | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--description` | string |  | help text shown for the path | — |
| `--position` | integer | `-1` | zero-based index to insert the path at (defaults to the end) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-path tenant --route create-user --start 1 --end 1
agnos add-path version --route api --start 0 --end 0 --trigger /v1
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
