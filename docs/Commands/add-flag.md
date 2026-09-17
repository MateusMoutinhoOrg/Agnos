# `add-flag`

Add a flag to a command's entries.yaml

```bash
agnos add-flag [--identifier <identifier>...] --command <command> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Appends one flag declaration to sandbox/internal/commands/<command>/entries.yaml and runs build so the command's new.go declares it. Without --identifier the flag answers to --<name>. Refuses a name or identifier the command already uses.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--identifier`, `-i` | string, repeatable |  | a cli identifier for the flag, e.g. --out or -o (repeatable; defaults to --<name>) |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that receives the field |
| `--type`, `-t` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description`, `-d` | string |  | help text shown for the field |
| `--example`, `-e` | string, repeatable |  | an usage example for the field (repeatable) |
| `--default` | string |  | the literal assigned when the field is absent (cannot be combined with --required) |
| `--required`, `-r` | boolean |  | fail with a usage error when the field is not provided (not for booleans or fields with --default) |
| `--array` | boolean |  | collect every occurrence into a []T field instead of a single value |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the field at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the flag name (the id the handler reads it back by, e.g. command.GetString("out-file")) |

```bash
agnos add-flag output --identifier --out --identifier -o --type string --required --command exec
agnos add-flag verbose --type boolean --description "print every step" --command exec
agnos add-flag retries --type int --min 0 --max 5 --default 1 --command exec
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
