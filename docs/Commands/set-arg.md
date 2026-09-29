# `set-arg`

Rewrite one arg of a command's command.yaml

```bash
agnos set-arg <Name> --command <command> [--rename <rename>] [--start <start>] [--end <end>] [--type <type>] [--required] [--default <default>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the arg to rewrite (its name or id) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that declares the arg | — |
| `--rename` | string |  | the name the arg takes on | — |
| `--start` | string |  | the first segment the arg reads | — |
| `--end` | string |  | the last segment the arg reads, -1 for the last one | — |
| `--type` | string |  | what one segment converts to: string, integer, number or uuid | — |
| `--required` | boolean |  | a command line matching the command without the arg is a usage error | — |
| `--default` | string |  | the literal bound when the arg is absent | — |
| `--trigger` | string |  | what the segments have to read as for the command to run | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix, text-prefix, suffix, regex or one-of | — |
| `--trigger-negate` | boolean |  | invert the trigger | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--description` | string |  | the one-line help text | — |
| `--clear` | string-array |  | a key to take off: trigger, trigger-negate, trigger-ignore-case, type, required, default or description (repeatable) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
