# Command Dispatch

## Description
Explains how a scaffolded project's command line is served: `agnos build` reads every `sandbox/internal/commands/<name>/entries.yaml` and generates `sandbox/internal/cli/climain.go`, whose `CliMain` picks the verb, and whose one `dispatch<Name>` function per command parses the remaining arguments into that command's `Entries` and calls its `CommandHandler`. Declaring a command is [ShapeCommands](/docs/ShapeCommands/doc.md); the keys the generator reads are [EntriesYaml](/docs/EntriesYaml/doc.md); the handler's contract is [WriteCommandHandler](/docs/WriteCommandHandler/doc.md).

---

## Three Files per Command

A command is a package under `sandbox/internal/commands/<name>/` holding exactly three files, and who writes each one is the whole contract:

| File | Written by | Holds |
|------|-----------|-------|
| `entries.yaml` | the project, through `add-command` / `add-flag` / `add-arg` / `set-command` | The declaration: identifiers, category, help, flags, args. |
| `entries.go` | `agnos build`, every time | `type Entries struct` — one typed field per flag, then per arg, in declaration order. |
| `handler.go` | the project, by hand | `func CommandHandler(deps *deps.Deps, entries *Entries) int`. |

There is no registry to append to. A directory with an `entries.yaml` is a command; `build` finds it by listing `sandbox/internal/commands/`.

---

## CliMain

`CliMain(deps, args)` is the single entry point `sandbox/api.Cli.CliMain` is bound to. It creates an argv parser through `deps.Argvdeps.New(args)`, drains the first positional as the verb, and switches on every identifier of every command:

```go
// sandbox/internal/cli/climain.go — generated
func CliMain(deps *deps.Deps, args []string) int {
	if len(args) == 0 {
		help.PrintGeneralHelp(deps)
		return ExitUsage // an empty command line is a usage error
	}

	verb := deps.Argvdeps.New(args)
	action, err := verb.GetNextStringArg()
	// …
	switch {
	case action == "greet":
		return dispatchGreet(deps, verb)
	case action == "help" || action == "--help":
		return dispatchHelp(deps, verb)
	case action == "version" || action == "--version":
		return dispatchVersion(deps, verb)
	}

	deps.Std.Error("unknown command %q — run '%s help' to see the available commands\n", action, binaryName())
	return ExitUsage
}
```

The exit-code constants are declared again inside the `cli` package, mirroring `sandbox/api`, so the generated layer depends on no contract package.

---

## The Dispatch Function

Each `dispatch<Name>` reads the command line in a fixed order, and every rejection exits `2` before the handler is reached:

1. **Flags, in declaration order.** A value flag is looked up by all its identifiers; when present, its raw text is converted to the field's type (`parseIntValue`, `parseFloatValue`, …), range-checked against `min`/`max`, and assigned. When absent, a `default` is assigned, a `required` field is a usage error, and anything else keeps the zero value. A boolean flag is a presence test. An `array` flag loops over `GetOptionsSize` occurrences.
2. **`quiet`.** If the command declares a boolean flag named `quiet` and it was given, `silenceLogs(deps)` replaces `deps.Std.Log` with a no-op right after reading it.
3. **Unknown flags.** `checkUnknownFlags` reports the first still-unread argument starting with `-`. A typo such as `--pathh` is an error, not a silently applied default.
4. **Args, in declaration order.** Each `nextArgValue` drains the next unread positional; conversion and range checks are the same as for flags, and an `array` arg drains everything left.
5. **Leftovers.** `checkUnusedArgs` reports the first argument nothing asked for — `agnos build .` fails here.
6. **The handler.** `return greet.CommandHandler(deps, entries)`.

```go
// generated for the greet example of ScaffoldProject.md
func dispatchGreet(deps *deps.Deps, verb argvdeps.Parser) int {
	entries := &greet.Entries{}
	if verb.GetOptionsSize([]string{"--name", "-n"}) > 0 {
		raw, rawOk := optionValue(deps, verb, "name", []string{"--name", "-n"}, 0)
		// … parseStringValue, assign
	} else {
		deps.Std.Error("required flag 'name' not provided\n")
		return ExitUsage
	}
	if !checkUnknownFlags(deps, verb) {
		return ExitUsage
	}
	if raw, rawOk := nextArgValue(verb); rawOk {
		value, valueOk := parseIntValue(deps, "arg", "times", raw)
		// … assign, then range check
		if entries.Times < 1 {
			deps.Std.Error("arg 'times' must be >= 1\n")
			return ExitUsage
		}
	} else {
		entries.Times = 1 // the declared default
	}
	if !checkUnusedArgs(deps, verb) {
		return ExitUsage
	}
	return greet.CommandHandler(deps, entries)
}
```

Because every rejection happens here, a handler that starts is holding valid, typed, in-range input, and never returns `ExitUsage` itself.

---

## The Help Command

`help` is a command like any other — the same three files under `sandbox/internal/commands/help/`, dispatched by `dispatchHelp` off its own identifiers (`help`, `--help`). The only difference is who writes the files: `agnos build` writes all three.

- `entries.yaml` is rendered from `assets/templates/help_entries.yaml` **before** the commands are collected, so the declaration is in the transaction when the collector reads it — but only when the file does not already exist. It is created once and then editable like any other command's declaration.
- `entries.go` is regenerated like every command's.
- `handler.go` comes from the `cli` asset group and bakes the collected metadata of every command, help included, into its `helpCommands` table: identifiers, category, help line, long description, examples, flags and args with their types, defaults and requiredness. Hand edits to it are overwritten on the next build.

The one special case in `CliMain` is the empty command line, which calls `help.PrintGeneralHelp` directly and exits `2`. The usage prefix the screens print (`demo greet`) is the project name from `config.ProjectName`, lowercased — the binary name a user types.

`add-command help` and `remove-command help` are refused: the command belongs to the generator.

---

## What the Sandbox Never Touches

`climain.go` reads the command line through `deps.Argvdeps` and prints through `deps.Std`; it never sees `os.Args` or `os.Stdout`. The binary in `cmd/main/main.go` — itself generated — is the only place the process's real argument vector is named:

```go
deps := agnosadapter.New()
lib := agnoslib.New(&deps)
os.Exit(lib.Cli.CliMain(os.Args[1:]))
```

So the whole interface runs against a fixed vector and a buffer as easily as against a terminal — see [ComposeDeps](/docs/ComposeDeps/doc.md).
