# `remove-dep`

Uninstalls one dep from the project

```bash
agnos remove-dep [--with-adapters] [--path <path>] [--quiet] <dep>
```

Removes every adapter whose declaration names the dep, its require and its enrollment in every available, then the contract itself, then calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--with-adapters` | boolean |  | also remove every adapter that fills the dep (without it, a dep with adapters installed is refused) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the dep to remove from the project |

```bash
agnos remove-dep embeddeps
agnos remove-dep embeddeps --path ./my-project
```

Deps System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
