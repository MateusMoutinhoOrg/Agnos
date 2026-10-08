# `set-command`

Rewrite the command-level keys of a command's command.yaml

```bash
agnos set-command <Name> [--summary <summary>] [--category <category>] [--description <description>] [--identifier <identifier>...] [--example <example>...] [--hidden] [--visible] [--priority <priority>] [--before <before>] [--after <after>] [--segments <segments>] [--strict] [--loose] [--clear <clear>...] [--help] [--path <path>] [--quiet]
```

Overwrites summary, category, description, priority, segments, strict and hidden on one command, appends further verbs (--identifier) and examples, and runs build. Empty options leave the current value alone; --before and --after place the command one rung from another; --clear takes a key off.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the command to update (a verb or its package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--summary` | string |  | new one-line help text | — |
| `--category` | string |  | new category the command is grouped under in help output | — |
| `--description` | string |  | new long description shown by help <command> | — |
| `--identifier`, `-i` | string-array |  | an extra verb the command answers to (repeatable) | — |
| `--example`, `-e` | string-array |  | an extra usage example (repeatable) | — |
| `--hidden` | boolean |  | hide the command from help listings | — |
| `--visible` | boolean |  | show the command in help listings again | — |
| `--priority` | integer | `-1` | the rung this command runs on, lowest first | — |
| `--before` | string |  | land one rung below the command named (excludes --priority) | — |
| `--after` | string |  | land one rung above the command named (excludes --priority) | — |
| `--segments` | integer | `-1` | how many segments the command line has to have for the command to run | — |
| `--strict` | boolean |  | every token of the line has to be read by the command or a middleware in front of it | — |
| `--loose` | boolean |  | turn the command into a middleware: tokens it does not read are the next command's to declare | — |
| `--clear` | string-array |  | a key to take off: segments, or examples before --example lists them again (repeatable) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos set-command exec --identifier run --example 'exec file.txt'
agnos set-command exec --priority 50
agnos set-command logger --loose
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
