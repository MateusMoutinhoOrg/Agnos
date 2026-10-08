# `help`

Display help for a command

```bash
agnos help [Name…] [--help] [--path <path>] [--quiet]
```

When called without arguments, lists every available command grouped by category. When called with a command name, shows detailed usage, arguments, flags, and examples for that command.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, repeatable |  | The command to describe; omit it to list every command |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos help
agnos help start
```

Info · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
