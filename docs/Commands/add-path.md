# `add-path`

Add one slice of the request path to a route

```bash
agnos add-path --route <route> [--start <start>] [--end <end>] [--type <type>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--position <position>] [--path <path>] [--quiet] <id>
```

Inserts one entry into the route's paths and runs build so the route's new.go and entries.go pick it up. A path reads the request segments from --start to --end (both inclusive, -1 the last one) as '/' followed by them joined by '/', binds that text to Entries.<Id>, and — with --trigger — only lets the route run when the text matches.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the path |
| `--start` | string |  | the index of the first request segment the slice reads (defaults to 0) |
| `--end` | string |  | the index of the last request segment the slice reads, -1 for the last one (defaults to -1) |
| `--type` | string |  | what the slice converts to: string (the default), integer, number or uuid; anything but string reads one segment (--start equal to --end), and a segment that does not convert makes the route a non-match |
| `--trigger` | string |  | what the slice has to read as for the route to run; without it the path is a plain capture |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (segment by segment), text-prefix, suffix or regex — or starts-with, ends-with, exact, matches (defaults to equal) |
| `--trigger-negate` | boolean |  | invert the trigger: the route runs when the slice does not match it |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case |
| `--description` | string |  | help text shown for the path |
| `--position` | int | `-1` | zero-based index to insert the path at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `id` | string, required |  | the id of the path, the Entries field its slice binds to (normalized to an exported Go name) |

```bash
agnos add-path tenant --route create-user --start 1 --end 1
agnos add-path version --route api --start 0 --end 0 --trigger /v1
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
