# `build`

Build the project in a directory

```bash
agnos build [--path <path>] [--quiet] [--runtime <runtime>] [--unsafe]
```

Re-renders every generated file of the project in the given directory, then hands the result to the runtime named by --runtime ("go" resolves the module graph and compiles every package, "none" renders only). If no path is provided, the current directory is used.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--runtime` | string | `go` | the toolchain the rendered project is handed to: go (tidy + compile) or none |
| `--unsafe` | boolean |  | Skips the verify schema gate before building |

```bash
agnos build
agnos build --path ./my-project
agnos build -q
```

Core Commands · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
