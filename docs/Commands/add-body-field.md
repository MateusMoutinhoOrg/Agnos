# `add-body-field`

Declare a property of a route's body json- or form-schema

```bash
agnos add-body-field <Name> --route <route> [--type <type>] [--required] [--array] [--min <min>] [--max <max>] [--exclusive-min <exclusive-min>] [--exclusive-max <exclusive-max>] [--format <format>] [--pattern <pattern>] [--enum <enum>...] [--const <const>] [--nullable] [--min-items <min-items>] [--max-items <max-items>] [--unique-items] [--additional-properties] [--no-additional-properties] [--help] [--path <path>] [--quiet]
```

Declares one property of the route's body schema at a dotted path, creating the objects it passes through, and runs build so the Body struct and BodySchema pick it up: the json-schema of a json body, the form-schema of a form one, which stays flat (no object, no nullable). A route that declared no body becomes a json one here. Every keyword the schema subset supports has a flag; ReadBody answers 400 on the first violation, naming the field path.

| Arg | Type | Default | Description |
| --- | --- | --- | --- |
| `Name` | string, required |  | the dotted path of the property (address.city); the objects it passes through are created as needed |

| Flag | Type | Default | Description | From |
| --- | --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the property | — |
| `--type` | string | `string` | the value type: string, boolean, int, float or object (defaults to string) | — |
| `--required` | boolean |  | list the property in its parent object's required set | — |
| `--array` | boolean |  | declare an array of the type instead of a single value | — |
| `--min` | string |  | minimum for a number, minLength for a string | — |
| `--max` | string |  | maximum for a number, maxLength for a string | — |
| `--exclusive-min` | string |  | exclusiveMinimum for a number property | — |
| `--exclusive-max` | string |  | exclusiveMaximum for a number property | — |
| `--format` | string |  | json-schema format for a string property: email, uuid, date-time or uri | — |
| `--pattern` | string |  | regular expression a string property must match | — |
| `--enum` | string-array |  | an accepted value of the property (repeatable; declares the enum set) | — |
| `--const` | string |  | the single value the property must carry | — |
| `--nullable` | boolean |  | accept null as well as the declared type | — |
| `--min-items` | string |  | shortest accepted array (--array only) | — |
| `--max-items` | string |  | longest accepted array (--array only) | — |
| `--unique-items` | boolean |  | refuse an array holding the same value twice (--array only) | — |
| `--additional-properties` | boolean |  | accept undeclared keys inside an object property | — |
| `--no-additional-properties` | boolean |  | refuse undeclared keys inside an object property | — |
| `--help`, `-h` | boolean |  | Print the help of the command this command line is for | [help-flag](help-flag.md) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) | [project](project.md) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output | [project](project.md) |

| Runs in front of it | When |
| --- | --- |
| [`help-flag`](help-flag.md) | always |
| [`project`](project.md) | always |

```bash
agnos add-body-field email --route create-user --format email --max 254 --required
agnos add-body-field address.city --route create-user --required
agnos add-body-field role --route create-user --enum admin --enum member
agnos add-body-field tags --route create-user --array --unique-items --max-items 10
```

Server System · [every command](doc.md) · [CommandYaml](../CommandYaml/doc.md)
