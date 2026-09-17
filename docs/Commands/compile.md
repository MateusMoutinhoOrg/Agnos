# `compile`

Cross-compile the project's binaries into release/

```bash
agnos compile --target <target>... [--path <path>] [--quiet]
```

Runs build over the project and then cross-compiles its ./cmd/main entrypoint once per --target into release/, with CGO disabled. Repeat --target for several targets, or pass --target all to build every one. Targets and their outputs: linux86 -> linux86.out, linuxarm64 -> linuxarm64.out, linuxi32 -> linuxi32.out, mac86 -> mac86.bin, macarm64 -> macarm64.bin, windows86 -> windows86.exe, windowsi32 -> windowsi32.exe.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--target`, `-t` | string, repeatable, required |  | a target to cross-compile (repeatable); one of linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32, or all |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos compile --target linux86
agnos compile --target linux86 --target macarm64
agnos compile --target all
```

Core Commands · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
