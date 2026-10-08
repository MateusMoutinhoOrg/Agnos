# `add-arg`

Add an arg to a command's command.yaml

```bash
agnos add-arg <Name> --command <command> [--start <start>] [--end <end>] [--type <type>] [--description <description>] [--default <default>] [--required] [--position <position>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--help] [--path <path>] [--quiet]
```

Inserts one arg — the segments --start to --end of the command line — into sandbox/internal/commands/<command>/command.yaml (at --position, else at the end) and runs build. An arg given no --start reads the first segment no arg reads yet; with --trigger it is part of what the command matches on.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the arg name; its exported Go form is the Input field the handler reads (file-name -> input.FileName) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that receives the arg | — |
| `--start` | string |  | the first segment the arg reads (defaults to the first no arg reads yet) | — |
| `--end` | string |  | the last segment the arg reads, -1 for the last one (defaults to --start) | — |
| `--type`, `-t` | string | `string` | what one segment converts to: string (the default), integer, number or uuid; anything but string reads one segment, and one that does not convert makes the command a non-match | — |
| `--description`, `-d` | string |  | the one-line help text of the arg | — |
| `--default` | string |  | the literal bound when the arg is absent (cannot be combined with --required) | — |
| `--required`, `-r` | boolean |  | a command line matching the command without the arg is a usage error | — |
| `--position` | integer | `-1` | zero-based index to insert the arg at (defaults to the end) | — |
| `--trigger` | string |  | what the segments, joined by a space, have to read as for the command to run; without it the arg is a plain capture | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix (word by word), text-prefix, suffix, regex or one-of (defaults to equal) | — |
| `--trigger-negate` | boolean |  | invert the trigger: the command runs when the segments do not match it | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-arg file --type string --required --description "the file to process" --command exec
agnos add-arg count --type integer --default 1 --position 0 --command exec
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
