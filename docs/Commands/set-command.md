# `set-command`

Update the command-level keys of a command's entries.yaml

```bash
agnos set-command [--help <help>] [--category <category>] [--long-description <long-description>] [--identifier <identifier>...] [--example <example>...] [--hidden] [--visible] [--path <path>] [--quiet] <name>
```

Rewrites help, category, long-description and hidden in sandbox/internal/commands/<name>/entries.yaml, and appends extra identifiers / examples, then runs build so help output is regenerated. Keys not passed are left untouched.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--help` | string |  | new one-line help text |
| `--category` | string |  | new category the command is grouped under in help output |
| `--long-description` | string |  | new long description shown by help <command> |
| `--identifier`, `-i` | string, repeatable |  | an extra verb the command answers to (repeatable) |
| `--example`, `-e` | string, repeatable |  | an extra usage example (repeatable) |
| `--hidden` | boolean |  | hide the command from help listings |
| `--visible` | boolean |  | show the command in help listings again |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the command to update (identifier or package name) |

```bash
agnos set-command exec --help "run the thing" --category Core
agnos set-command exec --identifier run --example "exec file.txt"
agnos set-command exec --hidden
```

Cli System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
