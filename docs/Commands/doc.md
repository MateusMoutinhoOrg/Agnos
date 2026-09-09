# Commands

`agnos <command> [flags] [args]`. `agnos help <command>` prints
the same for one command; an empty command line prints the general help and exits 2.

Hidden commands are not listed. Flags may appear anywhere on the command line; positionals
bind in order after them. A `repeatable` field is given once per value. Every section below is
rendered from that command's `entries.yaml` ([EntriesYaml](../EntriesYaml/doc.md)) on each
build.

## Cli System

### `add-arg`

Add a positional arg to a command's entries.yaml

```bash
agnos add-arg --command <command> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Inserts one positional arg declaration into sandbox/internal/commands/<command>/entries.yaml (at --position, else at the end) and runs build so entries.go and the dispatch layer are regenerated. Positional args bind by order; an array arg must stay last.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that receives the field |
| `--type`, `-t` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description`, `-d` | string |  | help text shown for the field |
| `--example`, `-e` | string, repeatable |  | an usage example for the field (repeatable) |
| `--default` | string |  | the literal assigned when the field is absent (cannot be combined with --required) |
| `--required`, `-r` | boolean |  | fail with a usage error when the field is not provided (not for booleans or fields with --default) |
| `--array` | boolean |  | collect every occurrence into a []T field instead of a single value |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the field at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the arg name (becomes the generated struct field) |

```bash
agnos add-arg file --type string --required --description "the file to process" --command exec
agnos add-arg count --type int --min 1 --position 0 --command exec
```

### `add-command`

Scaffold a new command package in the project

```bash
agnos add-command --help <help> --category <category> [--path <path>] [--quiet] <name>
```

Creates sandbox/internal/commands/<name>/ with a hand-written entries.yaml and a stub handler.go, then runs build so entries.go and the dispatch layer are generated for it. Refuses to overwrite an existing command.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--help` | string, required |  | one-line help text for the new command |
| `--category` | string, required |  | the category the new command is grouped under in help output |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new command (e.g. my-feature) |

```bash
agnos add-command my-feature
agnos add-command my-feature --path ./my-project
```

### `add-flag`

Add a flag to a command's entries.yaml

```bash
agnos add-flag [--identifier <identifier>...] --command <command> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Appends one flag declaration to sandbox/internal/commands/<command>/entries.yaml and runs build so entries.go and the dispatch layer are regenerated. Without --identifier the flag answers to --<name>. Refuses a name or identifier the command already uses.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--identifier`, `-i` | string, repeatable |  | a cli identifier for the flag, e.g. --out or -o (repeatable; defaults to --<name>) |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that receives the field |
| `--type`, `-t` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description`, `-d` | string |  | help text shown for the field |
| `--example`, `-e` | string, repeatable |  | an usage example for the field (repeatable) |
| `--default` | string |  | the literal assigned when the field is absent (cannot be combined with --required) |
| `--required`, `-r` | boolean |  | fail with a usage error when the field is not provided (not for booleans or fields with --default) |
| `--array` | boolean |  | collect every occurrence into a []T field instead of a single value |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the field at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the flag name (becomes the generated struct field, e.g. out-file -> OutFile) |

```bash
agnos add-flag output --identifier --out --identifier -o --type string --required --command exec
agnos add-flag verbose --type boolean --description "print every step" --command exec
agnos add-flag retries --type int --min 0 --max 5 --default 1 --command exec
```

### `cli-init`

Initializes the CLI layer for the project

```bash
agnos cli-init [--path <path>] [--quiet]
```

Installs the std and argv deps the CLI layer depends on, renders the "cli" asset group into the project, and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos cli-init
agnos cli-init --path ./my-project
```

### `cli-purge`

Removes the CLI layer from the project

```bash
agnos cli-purge [--path <path>] [--quiet]
```

Removes every file the "cli" asset group installs and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos cli-purge
agnos cli-purge --path ./my-project
```

### `remove-arg`

Remove a positional arg from a command's entries.yaml

```bash
agnos remove-arg --command <command> [--path <path>] [--quiet] <name>
```

Drops one positional arg declaration from sandbox/internal/commands/<command>/entries.yaml and runs build so entries.go and the dispatch layer forget it. Later args shift up.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the arg |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the arg name |

```bash
agnos remove-arg file --command exec
```

### `remove-command`

Delete a command package from the project

```bash
agnos remove-command [--path <path>] [--quiet] <name>
```

Deletes sandbox/internal/commands/<name>/ (entries.yaml, entries.go, handler.go and anything else inside) and runs build so climain.go and help stop dispatching to it. The generated help command cannot be removed.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the command to delete (identifier or package name) |

```bash
agnos remove-command my-feature
agnos remove-command my-feature --path ./my-project
```

### `remove-flag`

Remove a flag from a command's entries.yaml

```bash
agnos remove-flag --command <command> [--path <path>] [--quiet] <name>
```

Drops one flag declaration (matched by its name or by one of its identifiers) from sandbox/internal/commands/<command>/entries.yaml and runs build so entries.go and the dispatch layer forget it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--command`, `-c` | string, required |  | the command (identifier or package name) that owns the flag |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the flag name (or one of its identifiers, e.g. --out) |

```bash
agnos remove-flag output --command exec
agnos remove-flag --out --command exec
```

### `set-command`

Update the command-level keys of a command's entries.yaml

```bash
agnos set-command [--help <help>] [--category <category>] [--long-description <long-description>] [--identifier <identifier>...] [--example <example>...] [--hidden] [--visible] [--path <path>] [--quiet] <name>
```

Rewrites help, category, long-description and hidden in sandbox/internal/commands/<name>/entries.yaml, and appends extra identifiers / examples, then runs build so help output is regenerated. Keys not passed are left untouched.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--help` | string |  | new one-line help text |
| `--category` | string |  | new category the command is grouped under in help output |
| `--long-description` | string |  | new long description shown by help <command> |
| `--identifier`, `-i` | string, repeatable |  | an extra verb the command answers to (repeatable) |
| `--example`, `-e` | string, repeatable |  | an extra usage example (repeatable) |
| `--hidden` | boolean |  | hide the command from help listings |
| `--visible` | boolean |  | show the command in help listings again |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the command to update (identifier or package name) |

```bash
agnos set-command exec --help "run the thing" --category Core
agnos set-command exec --identifier run --example "exec file.txt"
agnos set-command exec --hidden
```

## Server System

### `add-body-field`

Declare a property of a route's body json-schema

```bash
agnos add-body-field --route <route> [--type <type>] [--required] [--array] [--min <min>] [--max <max>] [--exclusive-min <exclusive-min>] [--exclusive-max <exclusive-max>] [--format <format>] [--pattern <pattern>] [--enum <enum>...] [--const <const>] [--nullable] [--min-items <min-items>] [--max-items <max-items>] [--unique-items] [--additional-properties] [--no-additional-properties] [--path <path>] [--quiet] <name>
```

Declares one property of the route's body json-schema at a dotted path, creating the objects it passes through, and runs build so the Body struct and EntriesSchema pick it up. A route that declared no body becomes a json one here. Every keyword the schema subset supports has a flag; ReadBody answers 400 on the first violation, naming the field path.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the property |
| `--type` | string | `string` | the value type: string, boolean, int, float or object (defaults to string) |
| `--required` | boolean |  | list the property in its parent object's required set |
| `--array` | boolean |  | declare an array of the type instead of a single value |
| `--min` | string |  | minimum for a number, minLength for a string |
| `--max` | string |  | maximum for a number, maxLength for a string |
| `--exclusive-min` | string |  | exclusiveMinimum for a number property |
| `--exclusive-max` | string |  | exclusiveMaximum for a number property |
| `--format` | string |  | json-schema format for a string property: email, uuid, date-time or uri |
| `--pattern` | string |  | regular expression a string property must match |
| `--enum` | string, repeatable |  | an accepted value of the property (repeatable; declares the enum set) |
| `--const` | string |  | the single value the property must carry |
| `--nullable` | boolean |  | accept null as well as the declared type |
| `--min-items` | string |  | shortest accepted array (--array only) |
| `--max-items` | string |  | longest accepted array (--array only) |
| `--unique-items` | boolean |  | refuse an array holding the same value twice (--array only) |
| `--additional-properties` | boolean |  | accept undeclared keys inside an object property |
| `--no-additional-properties` | boolean |  | refuse undeclared keys inside an object property |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the dotted path of the property (address.city); the objects it passes through are created as needed |

```bash
agnos add-body-field email --route create-user --format email --max 254 --required
agnos add-body-field address.city --route create-user --required
agnos add-body-field role --route create-user --enum admin --enum member
agnos add-body-field tags --route create-user --array --unique-items --max-items 10
```

### `add-header`

Declare a request header on a route

```bash
agnos add-header --route <route> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Declares one request header on a route and runs build so entries.go and the dispatch arm pick it up. The name is the external spelling and is matched without regard to case; the dispatch answers 400 for a missing --required header or one outside --min/--max, before the handler runs.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the header |
| `--type` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description` | string |  | help text shown for the header |
| `--example` | string, repeatable |  | an usage example for the header (repeatable) |
| `--default` | string |  | the literal assigned when the header is absent (cannot be combined with --required) |
| `--required` | boolean |  | answer 400 when the header is not provided (not for booleans or headers with --default) |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the header at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the header name, matched without regard to case |

```bash
agnos add-header authorization --route create-user --required
agnos add-header x-retries --route create-user --type int --default 1 --max 5
```

### `add-param`

Declare a query parameter on a route

```bash
agnos add-param --route <route> [--type <type>] [--description <description>] [--example <example>...] [--default <default>] [--required] [--array] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] <name>
```

Declares one query-string parameter on a route and runs build so entries.go and the dispatch arm pick it up. --array is accepted here and nowhere else: every occurrence of the key is collected into a []T field.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the parameter |
| `--type` | string | `string` | the value type: string, boolean, int or float (defaults to string) |
| `--description` | string |  | help text shown for the parameter |
| `--example` | string, repeatable |  | an usage example for the parameter (repeatable) |
| `--default` | string |  | the literal assigned when the parameter is absent (cannot be combined with --required) |
| `--required` | boolean |  | answer 400 when the parameter is not provided (not for booleans or parameters with --default) |
| `--array` | boolean |  | collect every occurrence into a []T field instead of a single value |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the parameter at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query key |

```bash
agnos add-param page --route list-users --type int --default 1 --min 1
agnos add-param tag --route list-users --array
```

### `add-route`

Declare a new http route

```bash
agnos add-route [--trigger <trigger>] [--method <method>] --help <help> --category <category> [--path <path>] [--quiet] <name>
```

Writes sandbox/internal/routes/<name>/route.yaml and a stub handler.go, then runs build so entries.go and the dispatch arm are generated. The trigger is normalized to start with /, and defaults to /<name>.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--trigger` | string |  | the first literal segment of the path, always starting with / (defaults to /<name>) |
| `--method`, `-m` | string | `GET` | the http method the route answers: GET, POST, PUT, PATCH, DELETE, HEAD or OPTIONS |
| `--help` | string, required |  | one-line description of the route |
| `--category` | string, required |  | the heading the route is listed under in docs/Routes |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route name (becomes the directory sandbox/internal/routes/<name> and its Go package) |

```bash
agnos add-route create-user --trigger /users --method POST --help "Create a user" --category Users
```

### `add-segment`

Add a segment to a route's path

```bash
agnos add-segment --route <route> [--identifier <identifier>] [--type <type>] [--description <description>] [--example <example>...] [--min <min>] [--max <max>] [--position <position>] [--path <path>] [--quiet] [<name>]
```

Appends one segment to the route's paths and runs build so entries.go and the dispatch arm pick it up. With --identifier the segment is a literal, normalized to start with /; with a name it is a capture, which is always required and becomes an Entries field already converted.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) that receives the segment |
| `--identifier` | string |  | declare a literal segment instead of a capture, always starting with / |
| `--type` | string | `string` | the value type of a captured segment: string, boolean, int or float |
| `--description` | string |  | help text shown for the captured segment |
| `--example` | string, repeatable |  | an usage example for the segment (repeatable) |
| `--min` | string |  | smallest accepted value (int/float only) |
| `--max` | string |  | largest accepted value (int/float only) |
| `--position` | int | `-1` | zero-based index to insert the segment at (defaults to the end) |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string |  | the captured segment's name (omitted when --identifier declares a literal segment) |

```bash
agnos add-segment --route create-user --identifier /users
agnos add-segment tenant --route create-user --description "the tenant the user belongs to"
agnos add-segment page --route list-users --type int --min 1
```

### `remove-body-field`

Delete one property from a route's body json-schema

```bash
agnos remove-body-field --route <route> [--path <path>] [--quiet] <name>
```

Drops one property of the body json-schema, named by the same dotted path add-body-field declared it with, and unlists it from its parent's required set. The build renders only: dropping a property may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the property is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the dotted path of the property to drop |

```bash
agnos remove-body-field address.city --route create-user
```

### `remove-header`

Delete one declared header from a route

```bash
agnos remove-header --route <route> [--path <path>] [--quiet] <name>
```

Drops one declared header from a route, the exact inverse of add-header. The build renders only: dropping a header may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the header is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the header to drop |

```bash
agnos remove-header authorization --route create-user
```

### `remove-param`

Delete one declared query parameter from a route

```bash
agnos remove-param --route <route> [--path <path>] [--quiet] <name>
```

Drops one declared query parameter from a route, the exact inverse of add-param. The build renders only: dropping a parameter may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the parameter is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the query parameter to drop |

```bash
agnos remove-param page --route list-users
```

### `remove-route`

Delete one declared route

```bash
agnos remove-route [--path <path>] [--quiet] <name>
```

Removes sandbox/internal/routes/<name>/ whole and re-renders the dispatch. The build renders only: dropping a route may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the route to delete (identifier or package name) |

```bash
agnos remove-route create-user
```

### `remove-segment`

Delete one segment from a route's path

```bash
agnos remove-segment --route <route> [--path <path>] [--quiet] <name>
```

Drops one segment from the route's paths, the exact inverse of add-segment: a capture by its name, a literal by the identifier it spells. The build renders only: dropping a segment may leave hand-written code referring to what is gone.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--route` | string, required |  | the route (identifier or package name) the segment is declared on |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the segment to drop: the captured segment's name, or the identifier of a literal one |

```bash
agnos remove-segment tenant --route create-user
agnos remove-segment /users --route create-user
```

### `server-init`

Add the http server layer to the project

```bash
agnos server-init [--path <path>] [--quiet]
```

Installs the deps the server layer needs, renders sandbox/internal/server, the routeio package and the built-in health route, and writes the start-server command. A project with no cli layer is given one first: a server needs a command that starts it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos server-init
agnos server-init --path ./my-project
```

### `server-purge`

Remove the http server layer and every route in it

```bash
agnos server-purge [--path <path>] [--quiet]
```

Drops sandbox/internal/{server,routes,routeio} and the start-server command, then re-renders. The cli layer and the installed deps are left in place.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos server-purge
```

### `set-body`

Rewrite the body keys of a route.yaml

```bash
agnos set-body [--type <type>] [--required] [--optional] [--max-bytes <max-bytes>] [--content-type <content-type>] [--drop-schema] [--path <path>] [--quiet] <route>
```

Overwrites the body keys of one route.yaml: how the body is read, whether it is required, the longest one accepted and the content-type the dispatch demands. Empty options leave the current value alone. Declaring a json-schema is add-body-field's job; --drop-schema deletes the one already declared.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--type` | string |  | how the body is read: none, raw, text or json |
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

### `set-route`

Rewrite the route-level keys of a route.yaml

```bash
agnos set-route [--method <method>] [--help <help>] [--category <category>] [--long-description <long-description>] [--hidden] [--visible] [--path <path>] [--quiet] [--example <example>...] <route>
```

Overwrites method, help, category, long-description, hidden and examples on one route. Empty options leave the current value alone; --example appends.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--method`, `-m` | string |  | the http method the route answers |
| `--help` | string |  | one-line description of the route |
| `--category` | string |  | the heading the route is listed under in docs/Routes |
| `--long-description` | string |  | the paragraph docs/Routes prints under the route |
| `--hidden` | boolean |  | drop the route from docs/Routes, still dispatched |
| `--visible` | boolean |  | list the route again in docs/Routes |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--example` | string, repeatable |  | an usage example for the route (repeatable) |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `route` | string, required |  | the route to edit (identifier or package name) |

```bash
agnos set-route create-user --method POST --example "curl -X POST localhost:8080/users"
```

## Examples

### `add-cli-example`

Scaffold a new example under examples/cli/

```bash
agnos add-cli-example [--path <path>] [--quiet] <name>
```

Creates examples/cli/<name>/ with an example.sh stub that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name, and refuses outright in a project with no cli.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new example (it becomes one directory under examples/cli/) |

```bash
agnos add-cli-example start
```

### `add-lib-example`

Scaffold a new example under examples/lib/

```bash
agnos add-lib-example [--path <path>] [--quiet] <name>
```

Creates examples/lib/<name>/ with an example.go stub (package main) that already runs and exits 0, then runs build so the example listing of the docs names it. The golden result.yaml is written by the first exec-test, never by hand. Refuses an existing name.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the name of the new example (it becomes one directory under examples/lib/) |

```bash
agnos add-lib-example start
```

### `exec-test`

Run the project's examples and check them against their goldens

```bash
agnos exec-test [--only <only>] [--update] [--path <path>] [--quiet]
```

Runs every example of examples/cli/ and examples/lib/ in alphabetical order, cli side first, each with its own directory as the working directory and the project's own cli in front of the PATH. Every run starts from a removed TestDir and a removed AssertDir, and what it produced - the merged output, the exit status and the sha256 of every file the example copied out of TestDir into AssertDir - is compared against the example's result.yaml, or written there when that golden does not exist yet. An example that copied nothing out fails: it asserted nothing. An example declared on both sides must leave the same tree and exit the same way: the cli is only a wrapper over the lib.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--only` | string |  | run a single example by name, both sides (defaults to every example) |
| `--update` | boolean |  | rewrite every golden result.yaml with what this run produced instead of comparing |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos exec-test
agnos exec-test --only start
agnos exec-test --update
```

### `remove-cli-example`

Delete an example from examples/cli/

```bash
agnos remove-cli-example [--path <path>] [--quiet] <name>
```

Deletes examples/cli/<name>/ whole - the example.sh, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the example directory under examples/cli/ |

```bash
agnos remove-cli-example start
```

### `remove-lib-example`

Delete an example from examples/lib/

```bash
agnos remove-lib-example [--path <path>] [--quiet] <name>
```

Deletes examples/lib/<name>/ whole - the example.go, the golden result.yaml and any TestDir or AssertDir the last run left behind - and runs build so the example listing of the docs is rewritten without it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the example directory under examples/lib/ |

```bash
agnos remove-lib-example start
```

### `update-test`

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

## Documentation

### `add-doc`

Scaffold a new doc directory under docs/

```bash
agnos add-doc [--theme <theme>...] --description <description> [--path <path>] [--quiet] <name>
```

Creates docs/<name>/ with a doc.md stub and the props.yaml declaring it, then runs build so README.md's index and the parent's Index.md list it. A first-level doc needs at least one --theme of themes.yaml; a nested name (docs/<Parent>/<Name>) creates a sub-doc, which takes no theme. Refuses to overwrite an existing doc.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--theme`, `-t` | string, repeatable |  | a theme id of themes.yaml the doc belongs to (repeatable; first-level docs only) |
| `--description`, `-d` | string, required |  | the one-line summary every index lists the doc with |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc) |

```bash
agnos add-doc HandleReports --theme development --description "How a report is written and regenerated"
agnos add-doc PublicApi/api.AddDoc --description "The AddDoc action of the sandbox api"
```

### `remove-doc`

Delete a doc directory from docs/

```bash
agnos remove-doc [--path <path>] [--quiet] <name>
```

Deletes docs/<name>/ (doc.md, props.yaml, its assets and every sub-doc nested under it) and runs build so the indexes that listed it are rewritten without it. A theme left with no docs simply stops rendering a section in README.md; it is not an error, so themes.yaml can keep it.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `name` | string, required |  | the doc directory under docs/, nested with / for a sub-doc (e.g. PublicApi/api.AddDoc) |

```bash
agnos remove-doc HandleReports
agnos remove-doc PublicApi/api.AddDoc --path ./my-project
```

## Core Commands

### `build`

Build the project in a directory

```bash
agnos build [--path <path>] [--quiet] [--runtime <runtime>] [--unsafe]
```

Re-renders every generated file of the project in the given directory, then hands the result to the runtime named by --runtime ("go" resolves the module graph and compiles every package, "none" renders only). If no path is provided, the current directory is used.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--runtime` | string | `go` | the toolchain the rendered project is handed to: go (tidy + compile) or none |
| `--unsafe` | boolean |  | Skips the verify schema gate before building |

```bash
agnos build
agnos build --path ./my-project
agnos build -q
```

### `compile`

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

### `local-install`

Builds the project and installs it locally

```bash
agnos local-install [--path <path>] [--quiet]
```

Runs build over the project, then compiles ./cmd/main into /usr/local/bin/<project-name> (~/.local/bin on Windows) so the binary is on PATH.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

### `publish`

Builds, compiles and publishes a release via gh

```bash
agnos publish [--path <path>] [--release-name <release_name>] [--draft] [--target <target>] [--publisher <publisher>]
```

Runs build, then compile (every target by default), and publishes every file of release/ as a gh release named --release-name, defaulting to the version in AgnosConfig/project.yaml.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path`, `-p` | string | `.` | The directory holding the project (defaults to the current directory) |
| `--release-name`, `-rn` | string |  | The name of the release |
| `--draft` | boolean |  | Create a draft release |
| `--target`, `-t` | string | `all` | The target to compile for (defaults to all) |
| `--publisher`, `-pub` | string | `gh` | The publisher to use (defaults to gh) |

### `start`

Initialize a new project in a directory

```bash
agnos start [--path <path>] --project-name <project-name> [--quiet] [--force] [--module <module>]
```

Scaffolds a new Agnos project in the given directory, creating the required configuration files and folder structure. If no path is provided, the current directory is used.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--project-name`, `-p` | string, required |  | the name of the project |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |
| `--force`, `-f` | boolean |  | Forces the creation of the project, overwriting existing files |
| `--module`, `-m` | string |  | the go module path written into go.mod (required when the target dir has no go.mod yet) |

```bash
agnos start -p my-project
agnos start -p my-project --path ./my-project-dir
agnos start -p my-project -q
```

### `verify`

Checks the project keeps the sandbox/adapter schema

```bash
agnos verify [--path <path>] [--runtime <runtime>] [--quiet]
```

Verifies the structural rules the harness depends on: sandbox/ imports stay inside sandbox/, sandbox/ holds only api, binds, deps, internal and new.go, sandbox/api and sandbox/deps import nothing external, every sandbox/binds file mirrors a sandbox/api file and declares only functions, and adapters/ holds only availables and libs. `agnos build` runs this as a gate unless --unsafe is passed.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--runtime` | string | `go` | the toolchain the project is handed to after the schema check: go (tidy + compile) or none |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos verify
```

## Dependencies

### `dep-install`

Installs an embedded dep into the project

```bash
agnos dep-install [--path <path>] [--quiet] <dep>
```

Renders every file under assets/deplist/<dep> into the project at the path it holds inside that dep, then calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the dep to install from assets/deplist |

```bash
agnos dep-install embeddeps
agnos dep-install embeddeps --path ./my-project
```

### `dep-list`

Lists the embedded deps available to install

```bash
agnos dep-list [--path <path>] [--quiet]
```

Lists the name of every dep under assets/deplist that dep-install can render into a project.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos dep-list
```

### `dep-remove`

Removes an embedded dep from the project

```bash
agnos dep-remove [--path <path>] [--quiet] <dep>
```

Removes every file that assets/deplist/<dep> installs into the project, then calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `dep` | string, required |  | the dep to remove from the project |

```bash
agnos dep-remove embeddeps
agnos dep-remove embeddeps --path ./my-project
```

## Dependency System

### `deps-init`

Initializes the dependency-injection subsystem for the project

```bash
agnos deps-init [--path <path>] [--quiet]
```

Creates the sandbox/deps and adapters directories and calls build. Run this once before using dep-install.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos deps-init
agnos deps-init --path ./my-project
```

### `deps-purge`

Removes the dependency-injection subsystem from the project

```bash
agnos deps-purge [--path <path>] [--quiet]
```

Removes the sandbox/deps and adapters directories and calls build.

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--path` | string | `.` | the dir holding the project (defaults to the current directory) |
| `--quiet`, `-q` | boolean |  | Quiets the cli output |

```bash
agnos deps-purge
agnos deps-purge --path ./my-project
```

## Info

### `help` — `--help`

Display help for a command

```bash
agnos help [<command>]
```

When called without arguments, lists every available command grouped by category. When called with a command name, shows detailed usage, arguments, flags, and examples for that command.

| Argument | Type | Default | Description |
| --- | --- | --- | --- |
| `command` | string |  | The command to describe; omit it to list every command |

```bash
agnos help
agnos help start
```

### `version` — `--version`

Print the installed version

```bash
agnos version
```

Prints the current version of the installed binary and exits.

```bash
agnos version
```

Output channels and exit codes are in [Rules](../Rules/doc.md#output-channels).
