# `build`

Build the project in a directory

```bash
agnos build [--runtime <runtime>] [--unsafe] [--help] [--path <path>] [--quiet]
```

Re-renders every generated file of the project in the given directory, then hands the result to the runtime named by --runtime ("go" resolves the module graph and compiles every package, "none" renders only). If no path is provided, the current directory is used.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--runtime` | string | `go` | the toolchain the rendered project is handed to: go (tidy + compile) or none | — |
| `--unsafe` | boolean |  | Skips the verify schema gate before building | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos build
agnos build --path ./my-project
agnos build -q
```

Core Commands · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
