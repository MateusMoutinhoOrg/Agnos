# `set-flag`

Rewrite one flag of a command's command.yaml

```bash
agnos set-flag <Name> --command <command> [--rename <rename>] [--key <key>...] [--type <type>] [--required] [--default <default>] [--min <min>] [--max <max>] [--enum <enum>...] [--pattern <pattern>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--description <description>] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the flag to rewrite (its name, its id or one of its keys) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that declares the flag | — |
| `--rename` | string |  | the name the flag takes on | — |
| `--key` | string-array |  | a spelling the flag answers to; the ones given replace the declared ones (repeatable) | — |
| `--type` | string |  | string, integer, number, boolean, string-array or integer-array | — |
| `--required` | boolean |  | a command line without the flag is a usage error | — |
| `--default` | string |  | the literal bound when the flag is absent | — |
| `--min` | string |  | smallest accepted value (integer/number only) | — |
| `--max` | string |  | largest accepted value (integer/number only) | — |
| `--enum` | string-array |  | a value the flag accepts; the ones given replace the declared ones (repeatable) | — |
| `--pattern` | string |  | a regular expression every value has to match | — |
| `--trigger` | string |  | what the value has to read as for the command to run | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix, text-prefix, suffix, regex or one-of | — |
| `--trigger-negate` | boolean |  | invert the trigger | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--description` | string |  | the one-line help text | — |
| `--clear` | string-array |  | a key to take off: keys, type, required, default, min, max, enum, pattern, trigger, trigger-negate, trigger-ignore-case or description (repeatable) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
