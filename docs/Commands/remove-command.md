# `remove-command`

Delete a command package from the project

```bash
agnos remove-command <Name> [--help] [--path <path>] [--quiet]
```

Deletes the command's directory, in whatever folder of sandbox/internal/commands it sits (command.yaml, new.go, input.go, handler.go and anything else in it), and every folder that leaves empty, then runs build so the dispatch stops running it. A directory holding another command is refused, and so are help, version and help-flag, which are generated.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the command to delete (a verb or its package name) |

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
agnos remove-command my-feature
agnos remove-command my-feature --path ./my-project
```

Cli · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
