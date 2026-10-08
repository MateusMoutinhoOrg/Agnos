# `add-flag`

Add a flag to a command's command.yaml

```bash
agnos add-flag <Name> [--key <key>...] --command <command> [--type <type>] [--description <description>] [--default <default>] [--required] [--min <min>] [--max <max>] [--position <position>] [--enum <enum>...] [--pattern <pattern>] [--trigger <trigger>] [--trigger-type <trigger-type>] [--trigger-negate] [--trigger-ignore-case] [--help] [--path <path>] [--quiet]
```

Appends one flag declaration to sandbox/internal/commands/<command>/command.yaml and runs build so the command's new.go declares it and its input.go carries the field. Without --key the flag answers to --<name>. Refuses a name or a key the command already uses.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the flag name; its exported Go form is the Input field the handler reads (out-file -> input.OutFile) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--key` | string-array |  | a spelling the flag answers to, e.g. --out or -o (repeatable; defaults to --<name>) | — |
| `--command`, `-c` | string, required |  | the command (a verb or its package name) that receives the flag | — |
| `--type`, `-t` | string | `string` | the value type: string, integer, number, boolean, string-array or integer-array (defaults to string) | — |
| `--description`, `-d` | string |  | the one-line help text of the flag | — |
| `--default` | string |  | the literal bound when the flag is absent (cannot be combined with --required) | — |
| `--required`, `-r` | boolean |  | a command line without the flag is a usage error (not for a boolean, nor with --default) | — |
| `--min` | string |  | smallest accepted value (integer/number only) | — |
| `--max` | string |  | largest accepted value (integer/number only) | — |
| `--position` | integer | `-1` | zero-based index to insert the flag at (defaults to the end) | — |
| `--enum` | string-array |  | a value the flag accepts, every other refused (repeatable) | — |
| `--pattern` | string |  | a regular expression every value has to match | — |
| `--trigger` | string |  | what the value has to read as for the command to run at all | — |
| `--trigger-type` | string |  | how the trigger is compared: equal, prefix, text-prefix, suffix, regex or one-of (defaults to equal) | — |
| `--trigger-negate` | boolean |  | invert the trigger: the command runs when the value does not match it | — |
| `--trigger-ignore-case` | boolean |  | compare the trigger without regard to case | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos add-flag output --key --out --key -o --required --command exec
agnos add-flag verbose --type boolean --description 'print every step' --command exec
agnos add-flag retries --type integer --min 0 --max 5 --default 1 --command exec
agnos add-flag tag --type string-array --enum a --enum b --command exec
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
