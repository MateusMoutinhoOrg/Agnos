# `update-test`

Rewrite one example's golden with what it produces now

```bash
agnos update-test [--path <path>] [--quiet] <name>
```

Runs one example by name, both sides, and writes what it produced over its result.yaml instead of comparing against it. Every write prints what it changes first - the paths that entered, left or changed sha, and the old output against the new one - so a golden is never rewritten unread. It is the normal way one golden is refreshed; exec-test --update rewrites the whole suite at once and hides the one that moved for a reason nobody meant.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the example to update, both sides |

```bash
agnos update-test start
agnos update-test add-command --path ./my-project
```

Examples · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
