# `set-body`

Rewrite the body keys of a route.yaml

```bash
agnos set-body [--type <type>] [--required] [--optional] [--max-bytes <max-bytes>] [--content-type <content-type>] [--drop-schema] [--path <path>] [--quiet] <route>
```

Overwrites the body keys of one route.yaml: how the body is read, whether it is required, the longest one accepted and the content-type the dispatch demands. Empty options leave the current value alone. Declaring a json-schema is add-body-field's job; --drop-schema deletes the one already declared.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--type` | string |  | how the body is read: none, raw, text, json or form (application/x-www-form-urlencoded, read into a map[string][]string) |
| `--required` | boolean |  | answer 400 when the body is absent or empty |
| `--optional` | boolean |  | accept an absent body again |
| `--max-bytes` | int | `-1` | the longest body accepted, in bytes; a longer one is answered 413 |
| `--content-type` | string |  | the only content-type accepted; a divergent one is answered 415 |
| `--drop-schema` | boolean |  | delete the declared json-schema, leaving the body unvalidated |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route to edit (identifier or package name) |

```bash
agnos set-body create-user --type json --required --max-bytes 2097152
agnos set-body upload-avatar --type raw --content-type application/octet-stream
agnos set-body ping --type none
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
