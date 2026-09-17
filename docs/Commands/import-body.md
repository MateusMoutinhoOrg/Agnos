# `import-body`

Infer a route's body json-schema from an example payload

```bash
agnos import-body [--json <json>] [--file <file>] [--required] [--replace] [--infer-format] [--path <path>] [--quiet] <route>
```

Reads one example payload and declares a body property for every key it carries, which is add-body-field run once per key. A route that declared no body becomes a json one here. The inference is a starting point: a type per key, the objects and lists around them, and — with --infer-format — the four formats a string may spell. A property already declared is never written over; --replace starts the schema over instead. Every bound after that is set-body-field's.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--json` | string |  | the example payload itself, as one json document |
| `--file` | string |  | a file holding the example payload, read instead of --json |
| `--required` | boolean |  | list every key the example carries in its object's required set |
| `--replace` | boolean |  | start the schema over instead of adding to the one declared |
| `--infer-format` | boolean |  | read an email, a uuid, a date-time or a uri back as the format it spells |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route (identifier or package name) that receives the schema |

```bash
agnos import-body create-user --file payload.json --required --infer-format
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
