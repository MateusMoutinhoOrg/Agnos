# Commands

`agnos <command> [flags] [args]`. `agnos help <command>` prints the same for one command.

Common flags: `--path <dir>` (default `.`, the only way to name the project; a positional dir is a usage error) and `--quiet`/`-q` (silences progress on stderr). Flags may appear anywhere; positionals bind in order after flags.

`Rebuild` column: the runtime the follow-up `build` runs with. `go` = `go mod tidy` + `go build`; `none` = render only.

## Project

| Command | Does | Rebuild |
|---|---|---|
| `start --project-name <n> [--module <m>] [--force]` | Writes `AgnosConfig/` and `go.mod` (`--module` required when none exists, `--force` overwrites), then builds | go |
| `build [--runtime go\|none] [--unsafe]` | `verify` (skipped by `--unsafe`), re-render generated files, persist, runtime | - |
| `verify [--runtime go\|none]` | Schema check, no writes, lists every violation, then runtime | - |
| `compile --target <t>... [--path]` | `build`, then `go build -o release/<file> ./cmd/main` per target with `CGO_ENABLED=0` | go |
| `local-install` | `build`, then `go build` into `/usr/local/bin/<name>` (Windows: `~/.local/bin`) | go |
| `publish [--release-name <v>] [--draft] [--target <t>] [--publisher gh]` | `build`, `compile` (default `all`), then `gh release create <v>` with every file of `release/`. `<v>` defaults to `version` in `project.yaml` | go |

Targets: `linux86 linuxarm64 linuxi32 mac86 macarm64 windows86 windowsi32` or `all`; outputs `linux86.out linuxarm64.out linuxi32.out mac86.bin macarm64.bin windows86.exe windowsi32.exe`.

## Deps

| Command | Does | Rebuild |
|---|---|---|
| `deps-init` | Creates `sandbox/deps/` + `adapters/` | go |
| `deps-purge` | Removes both whole | none |
| `dep-list` | Prints installable dep names (stdout) | - |
| `dep-install <dep>` | Renders `assets/deplist/<dep>/**` into the project, adds pinned `require` to `go.mod` if listed in `depsversion.yaml` | go |
| `dep-remove <dep>` | Deletes those files (and emptied dirs), strips the `require` | none |

## CLI layer

| Command | Does | Rebuild |
|---|---|---|
| `cli-init` | Installs `std` + `argvdeps`, renders the `cli` group (`cmd/main`, dispatch, `help`, `version`) | go |
| `cli-purge` | Removes the `cli` group files plus `sandbox/internal/{cli,commands}` whole. Leaves the two deps | none |
| `add-command <name> --help <t> --category <c>` | Writes `entries.yaml` + stub `handler.go`. Name normalized to `[a-z][a-z0-9-]*`; `help` refused | go |
| `remove-command <name>` | Deletes the command dir. `help` refused | none |
| `set-command <name> [--help] [--category] [--long-description] [--hidden\|--visible] [--identifier]... [--example]...` | Rewrites command-level keys; identifiers/examples append, deduplicated | go |
| `add-flag <name> --command <cmd> [--identifier <id>]... [--type string\|boolean\|int\|float] [--description] [--example]... [--default] [--required] [--array] [--min] [--max] [--position <i>]` | Appends a flag. Identifiers default to `--<name>`, type to `string` | go |
| `remove-flag <name-or-identifier> --command <cmd>` | Drops a flag | none |
| `add-arg <name> --command <cmd> [same field flags as add-flag]` | Inserts a positional at `--position` (else last). An `--array` arg must stay last | go |
| `remove-arg <name> --command <cmd>` | Drops a positional | none |

`--command`/`-c` takes the identifier or package name. `--required` is refused on a boolean or with `--default`. `--min`/`--max` apply to `int`/`float`. Keys are in [EntriesYaml](/docs/EntriesYaml/doc.md).

## Docs

| Command | Does | Rebuild |
|---|---|---|
| `add-doc <name> --description <t> [--theme <id>]...` | Writes `docs/<name>/{doc.md,props.yaml}`. `<name>` is `/`-nested for a sub-doc (parent must exist). `--theme` required on a first-level doc, refused on a sub-doc, must exist in `themes.yaml`. `Index` refused | go |
| `remove-doc <name>` | Deletes `docs/<name>/` with sub-docs and assets | none |

## Info

| Command | Does |
|---|---|
| `help [command]`, `--help` | General listing by category, or one command's usage. Empty command line prints help and exits 2 |
| `version`, `--version` | Prints `Version:<v>` from `project.yaml` |

## Output channels

| Channel | Stream | Carries | `--quiet` |
|---|---|---|---|
| `deps.Std.Printf` | stdout | The result (listings, version, help) | kept |
| `deps.Std.Log` | stderr | Progress | silenced |
| `deps.Std.Error` | stderr | Usage errors and failures | kept |

## Exit codes

| Code | Const | Meaning |
|---|---|---|
| 0 | `api.ExitOk` | Done |
| 1 | `api.ExitFailure` | Well-formed command failed (verify violation, compile error, unknown dep, existing command) |
| 2 | `api.ExitUsage` | Bad command line: unknown command/flag, leftover positional, missing required, bad or out-of-range number, unknown `--runtime` |
