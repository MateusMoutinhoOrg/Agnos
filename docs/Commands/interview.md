# `interview`

Guided mode: answer questions instead of typing commands

```bash
agnos interview [--path <path>] [--quiet]
```

The one screen of agnos made for a person rather than for a script: it reads the project you point it at and offers what that project can actually do next — create it, give it a command line, declare a command — one suggested step at a time. Areas you have not turned on stay off the menu until you turn them on. Every answer is a flag or an argument you could have typed, and the command line your answers add up to is shown before anything runs.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos interview
agnos interview --path ./my-project
```

Info · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
