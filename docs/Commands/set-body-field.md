# `set-body-field`

Rewrite one property of a route's body json-schema

```bash
agnos set-body-field --route <route> [--rename <rename>] [--type <type>] [--required] [--array] [--min <min>] [--max <max>] [--exclusive-min <exclusive-min>] [--exclusive-max <exclusive-max>] [--format <format>] [--pattern <pattern>] [--enum <enum>...] [--const <const>] [--nullable] [--min-items <min-items>] [--max-items <max-items>] [--unique-items] [--additional-properties] [--no-additional-properties] [--clear <clear>...] [--path <path>] [--quiet] <name>
```

Rewrites one property of the route's body json-schema in place and runs build. The keywords already declared are read back, the ones given are written over them, and the whole is built again by the constructor add-body-field uses — so the property a forgotten --max is added to is the property that was there. --clear takes a keyword off again, and a --type the old keywords cannot survive drops them, naming each one it dropped. A property that is not declared yet is add-body-field's.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the property is declared on |
| `--rename` | string |  | the key the property answers to from now on (it stays in the object it is declared in) |
| `--type` | string |  | the value type: string, boolean, int, float or object |
| `--required` | boolean |  | list the property in its parent object's required set |
| `--array` | boolean |  | declare an array of the type instead of a single value |
| `--min` | string |  | minimum for a number, minLength for a string |
| `--max` | string |  | maximum for a number, maxLength for a string |
| `--exclusive-min` | string |  | exclusiveMinimum for a number property |
| `--exclusive-max` | string |  | exclusiveMaximum for a number property |
| `--format` | string |  | json-schema format for a string property: email, uuid, date-time or uri |
| `--pattern` | string |  | regular expression a string property must match |
| `--enum` | string, repeatable |  | an accepted value of the property (repeatable; replaces the enum set) |
| `--const` | string |  | the single value the property must carry |
| `--nullable` | boolean |  | accept null as well as the declared type |
| `--min-items` | string |  | shortest accepted array (an array property only) |
| `--max-items` | string |  | longest accepted array (an array property only) |
| `--unique-items` | boolean |  | refuse an array holding the same value twice (an array property only) |
| `--additional-properties` | boolean |  | accept undeclared keys inside an object property |
| `--no-additional-properties` | boolean |  | refuse undeclared keys inside an object property |
| `--clear` | string, repeatable |  | a keyword to take off again: required, array, min, max, format, pattern, enum, const, nullable and the rest (repeatable) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the dotted path of the property to edit (address.city) |

```bash
agnos set-body-field age --route create-user --type int --min 0 --max 130
```

Server System · [every command](doc.md) · [EntriesYaml](../EntriesYaml/doc.md)
