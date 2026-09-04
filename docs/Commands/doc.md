# Commands

`agnos <command> [flags] [args]`. `agnos help <command>` prints the same for one
command; an empty command line prints the general help and exits 2.

Every section below is rendered from `sandbox/internal/commands/<name>/entries.yaml` on each
build: declare a command with `add-command`, its fields with `add-flag` / `add-arg`, and its
prose with `set-command`. Hidden commands are not listed. Flags may appear anywhere on the
command line; positionals bind in order after them. A `repeatable` field is given once per
value.

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
