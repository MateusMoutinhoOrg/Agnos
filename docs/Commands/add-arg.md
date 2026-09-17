# `add-arg`

Add a positional arg to a command's entries.yaml

```bash
agnos add-arg --command <command> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Inserts one positional arg declaration into sandbox/internal/commands/<command>/entries.yaml (at --position, else at the end) and runs build so the command's new.go declares it. Positional args bind by order; an array arg must stay last.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
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
| `name` | string, required |  | the arg name (the id the handler reads it back by) |

```bash
agnos add-arg file --type string --required --description "the file to process" --command exec
agnos add-arg count --type int --min 1 --position 0 --command exec
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
