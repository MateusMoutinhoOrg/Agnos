# `verify`

Check the project keeps the sandbox/adapter schema

```bash
agnos verify [--runtime <runtime>] [--help] [--path <path>] [--quiet]
```

Verifies the structural rules the harness depends on: sandbox/ imports stay inside sandbox/, sandbox/ holds only api, deps, internal and new.go, sandbox/api imports nothing but sandbox/deps and sandbox/deps imports nothing external, every sandbox/api file has the sandbox/internal/<x>/new.go that builds it, and adapters/ holds only bindings and libs. `agnos build` runs this as a gate unless --unsafe is passed.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--runtime` | string | `go` | the toolchain the project is handed to after the schema check: go (tidy + compile) or none | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos verify
```

Core · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
