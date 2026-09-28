# `compile`

Cross-compile the project's binaries into release/

```bash
agnos compile --target <target>... [--help] [--path <path>] [--quiet]
```

Runs build over the project and then cross-compiles its ./cmd/main entrypoint once per --target into release/, with CGO disabled. Repeat --target for several targets, or pass --target all to build every one. Targets and their outputs: linux86 -> linux86.out, linuxarm64 -> linuxarm64.out, linuxi32 -> linuxi32.out, mac86 -> mac86.bin, macarm64 -> macarm64.bin, windows86 -> windows86.exe, windowsi32 -> windowsi32.exe.

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--target`, `-t` | string-array, required |  | a target to cross-compile (repeatable); one of linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32, or all | — |
| `--help` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos compile --target linux86
agnos compile --target linux86 --target macarm64
agnos compile --target all
```

Core Commands · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
