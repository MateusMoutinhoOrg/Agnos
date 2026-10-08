# `set-body`

Rewrite the body keys of a route.yaml

```bash
agnos set-body <Name> [--type <type>] [--required] [--optional] [--max-bytes <max-bytes>] [--content-type <content-type>] [--drop-schema] [--help] [--path <path>] [--quiet]
```

Overwrites the body keys of one route.yaml: how the body is read, whether it is required, the longest one accepted and the content-type the dispatch demands. Empty options leave the current value alone. Declaring a schema is add-body-field's job; turning a json body into a form one (or back) carries its schema along when it is flat, and --drop-schema deletes the one already declared.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the route to edit (identifier or package name) |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--type` | string |  | how the body is read: none, raw, text, json or form (application/x-www-form-urlencoded: a map[string][]string, or the Body struct its form-schema describes) | — |
| `--required` | boolean |  | answer 400 when the body is absent or empty | — |
| `--optional` | boolean |  | accept an absent body again | — |
| `--max-bytes` | integer | `-1` | the longest body accepted, in bytes; a longer one is answered 413 | — |
| `--content-type` | string |  | the only content-type accepted; a divergent one is answered 415 | — |
| `--drop-schema` | boolean |  | delete the declared json- or form-schema, leaving the body unvalidated | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project-flags](project-flags.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project-flags](project-flags.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project-flags`](project-flags.md) | always |

```bash
agnos set-body create-user --type json --required --max-bytes 2097152
agnos set-body upload-avatar --type raw --content-type application/octet-stream
agnos set-body ping --type none
agnos set-body login --type form
```

Server · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
