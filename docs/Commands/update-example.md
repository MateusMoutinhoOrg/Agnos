# `update-example`

Rewrite one example's golden with what it produces now

```bash
agnos update-example <Name> [--help] [--path <path>] [--quiet]
```

Runs one example by name, both sides, and writes what it produced over its result.yaml instead of comparing against it. Every write prints what it changes first - the paths that entered, left or changed sha, and the old output against the new one - so a golden is never rewritten unread. It is the normal way one golden is refreshed; run-examples --update rewrites the whole suite at once and hides the one that moved for a reason nobody meant.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the example to update, both sides |

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
agnos update-example start
agnos update-example add-command --path ./my-project
```

Examples · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
