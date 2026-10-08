# `remove-dep`

Uninstall one dep from the project

```bash
agnos remove-dep <Name> [--with-adapters] [--help] [--path <path>] [--quiet]
```

Removes every adapter whose declaration names the dep, its require and its enrollment in every binding, then the contract itself, then calls build.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the dep to remove from the project |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--with-adapters` | boolean |  | also remove every adapter that fills the dep (without it, a dep with adapters installed is refused) | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos remove-dep embeddeps
agnos remove-dep embeddeps --path ./my-project
```

Deps · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
