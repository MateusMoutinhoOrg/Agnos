# `rename-command`

Rename a command: its package, and the verb it answered to by its name

```bash
agnos rename-command <Target> <Name> [--dir <dir>] [--help] [--path <path>] [--quiet]
```

Moves the command's directory to <name>/ in the folder it sits in — or in the one --dir names, / for the top — rewriting the package clause of every hand-written Go file and pointing the verb at the new name when it was the old one, drops every folder that leaves empty, and runs build so new.go, entries.go and the dispatch follow it. <name> may be the current one when only the folder changes. help, version and help-flag are generated and refused.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Target` | string, required |  | the command to rename (its name or a verb it answers to) |
| `Name` | string, required |  | the name it takes on |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--dir` | string |  | the folder under sandbox/internal/commands the command moves to, / for the top (defaults to the folder it sits in) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

Cli System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
