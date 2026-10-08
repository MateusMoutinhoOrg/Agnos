# `interview`

Answer questions instead of typing commands, in guided mode

```bash
agnos interview [--help] [--path <path>] [--quiet]
```

The one screen of agnos made for a person rather than for a script: it reads the project you point it at and offers what that project can actually do next — create it, give it a command line, declare a command — one suggested step at a time. Areas you have not turned on stay off the menu until you turn them on. Every answer is a flag or an argument you could have typed, and the command line your answers add up to is shown before anything runs.

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
agnos interview
agnos interview --path ./my-project
```

Info · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
