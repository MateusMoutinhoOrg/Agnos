package cli

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/argvdeps"
	"{{.Module}}/sandbox/internal/config"
)

// Exit codes. Kept here (not in sandbox/api) so the cli layer has no
// dependency on the contract package. They mirror sandbox/api: 0 success,
// 1 a well-formed command that failed, 2 a command line that was wrong.
const (
	ExitOk      = 0
	ExitFailure = 1
	ExitUsage   = 2
)

// The declared types a flag or an arg may carry, as entries.yaml spells them.
const (
	typeBoolean = "boolean"
	typeInt     = "int"
	typeFloat   = "float"
)

// subjectFlag and subjectArg are how a usage error names what it is about.
const (
	subjectFlag = "flag"
	subjectArg  = "arg"
)

// quietId is the flag every command that wants a silent run declares. The
// dispatch does the silencing once, so no handler has to.
const quietId = "quiet"

// helpId is the command the empty command line falls back to.
const helpId = "help"

// CliMain reads the verb, finds the command of Cli.Commands that answers to
// it, binds the rest of the command line to a copy of that command's declared
// flags and args, and calls its handler. Nothing here is generated per command:
// every command is one declaration built by its own NewCommand and collected by
// sandbox/internal/cli/new.go, so this dispatch is the same file in every
// project.
// `help` is reached through that same path — it is a declared command whose
// files `{{.GeneratorName}} build` happens to write itself — and directly only
// for the empty command line below.
func CliMain(sandbox *api.Sandbox, args []string) int {

	if len(args) == 0 {
		return runHelp(sandbox)
	}

	parser := sandbox.Deps.Argvdeps.New(args)

	action, err := parser.GetNextStringArg()
	if err != nil {
		return runHelp(sandbox)
	}

	for index := range sandbox.Cli.Commands {
		declared := &sandbox.Cli.Commands[index]
		if answersTo(declared, action) {
			return runCommand(sandbox, api.BindCommand(declared), parser)
		}
	}

	sandbox.Deps.Std.Error("unknown command %q — run '%s help' to see the available commands\n", action, binaryName(sandbox))
	return ExitUsage
}

// answersTo reports whether a command is spelled by this verb, its aliases
// included.
func answersTo(command *api.Command, action string) bool {
	for _, identifier := range command.Identifiers {
		if identifier == action {
			return true
		}
	}
	return false
}

// runHelp prints the general help screen by running the help command with
// nothing bound to it, and reports the empty command line as a usage error.
func runHelp(sandbox *api.Sandbox) int {
	for index := range sandbox.Cli.Commands {
		declared := &sandbox.Cli.Commands[index]
		if answersTo(declared, helpId) && declared.Handler != nil {
			command := api.BindCommand(declared)
			command.Handler(command)
			break
		}
	}
	return ExitUsage
}

// runCommand binds one command line onto one bound copy of a command's
// declaration and hands the result to its handler. The order is the order a
// user reads a command line in: the flags, then the leftovers that look like
// flags, then the positional args, then whatever is still unread.
func runCommand(sandbox *api.Sandbox, command *api.Command, parser argvdeps.Parser) int {
	for _, flag := range command.Flags {
		if !bindFlag(sandbox, command, flag, parser) {
			return ExitUsage
		}
	}

	if command.GetBool(quietId) {
		silenceLogs(sandbox)
	}

	if !checkUnknownFlags(sandbox, parser) {
		return ExitUsage
	}

	for _, arg := range command.Args {
		if !bindArg(sandbox, command, arg, parser) {
			return ExitUsage
		}
	}

	if !checkUnusedArgs(sandbox, parser) {
		return ExitUsage
	}

	return command.Handler(command)
}

// bindFlag reads one declared flag off the command line into command.Items.
func bindFlag(sandbox *api.Sandbox, command *api.Command, flag api.CommandFlag, parser argvdeps.Parser) bool {
	if flag.Type == typeBoolean {
		command.Items[flag.Id] = []any{parser.IsPresent(flag.Identifiers)}
		return true
	}

	occurrences := parser.GetOptionsSize(flag.Identifiers)

	if flag.Array {
		for occurrence := 0; occurrence < occurrences; occurrence++ {
			raw, rawOk := optionValue(sandbox, parser, flag.Id, flag.Identifiers, occurrence)
			if !rawOk {
				return false
			}
			value, valueOk := parseValue(sandbox, subjectFlag, flag.Id, flag.Type, raw)
			if !valueOk {
				return false
			}
			command.Items[flag.Id] = append(command.Items[flag.Id], value)
		}
		return bound(sandbox, command, flag.Id, subjectFlag, flag.Required)
	}

	if occurrences > 0 {
		raw, rawOk := optionValue(sandbox, parser, flag.Id, flag.Identifiers, 0)
		if !rawOk {
			return false
		}
		value, valueOk := parseValue(sandbox, subjectFlag, flag.Id, flag.Type, raw)
		if !valueOk {
			return false
		}
		if !inRange(sandbox, subjectFlag, flag.Id, flag.Type, value, flag.Min, flag.HasMin, flag.Max, flag.HasMax) {
			return false
		}
		command.Items[flag.Id] = []any{value}
		return true
	}

	if flag.Required {
		sandbox.Deps.Std.Error("required %s '%s' not provided\n", subjectFlag, flag.Id)
		return false
	}
	if flag.HasDefault {
		command.Items[flag.Id] = []any{defaultValue(sandbox, flag.Type, flag.Default)}
	}
	return true
}

// bindArg reads one declared positional argument off the command line into
// command.Items. An array arg drains every argument that is left.
func bindArg(sandbox *api.Sandbox, command *api.Command, arg api.CommandArg, parser argvdeps.Parser) bool {
	if arg.Array {
		for {
			raw, rawOk := nextArgValue(parser)
			if !rawOk {
				break
			}
			value, valueOk := parseValue(sandbox, subjectArg, arg.Id, arg.Type, raw)
			if !valueOk {
				return false
			}
			command.Items[arg.Id] = append(command.Items[arg.Id], value)
		}
		return bound(sandbox, command, arg.Id, subjectArg, arg.Required)
	}

	if raw, rawOk := nextArgValue(parser); rawOk {
		value, valueOk := parseValue(sandbox, subjectArg, arg.Id, arg.Type, raw)
		if !valueOk {
			return false
		}
		if !inRange(sandbox, subjectArg, arg.Id, arg.Type, value, arg.Min, arg.HasMin, arg.Max, arg.HasMax) {
			return false
		}
		command.Items[arg.Id] = []any{value}
		return true
	}

	if arg.Required {
		sandbox.Deps.Std.Error("required %s '%s' not provided\n", subjectArg, arg.Id)
		return false
	}
	if arg.HasDefault {
		command.Items[arg.Id] = []any{defaultValue(sandbox, arg.Type, arg.Default)}
	}
	return true
}

// bound reports a required array flag or arg that nothing was bound to.
func bound(sandbox *api.Sandbox, command *api.Command, id string, subject string, required bool) bool {
	if required && len(command.Items[id]) == 0 {
		sandbox.Deps.Std.Error("required %s '%s' not provided\n", subject, id)
		return false
	}
	return true
}

// ─── argv helpers ───────────────────────────────────────────────────────────
//
// Every value read above goes through these, so one spelling of every usage
// error is emitted for the whole CLI.

// binaryName is the executable's name as a user types it: the configured
// project name, lowercased.
func binaryName(sandbox *api.Sandbox) string {
	return sandbox.Deps.Stringsdeps.ToLower(config.ProjectName)
}

// silenceLogs turns off the progress channel for the rest of the process. It
// backs the --quiet flag: results (Printf) and errors (Error) still go out,
// only the "… started with path …" notices stop.
func silenceLogs(sandbox *api.Sandbox) {
	sandbox.Deps.Std.Log = func(format string, a ...any) (int, error) {
		return 0, nil
	}
}

// optionValue reads the occurrence-th value of a value flag, reporting a
// clean usage error when the flag was given with nothing after it.
func optionValue(sandbox *api.Sandbox, parser argvdeps.Parser, name string, flags []string, occurrence int) (string, bool) {
	raw, err := parser.GetStringOption(flags, occurrence)
	if err != nil {
		sandbox.Deps.Std.Error("flag '%s': expected a value after %s\n", name, flags[0])
		return "", false
	}
	return raw, true
}

// nextArgValue drains the next positional argument, reporting whether one was
// left. It never fails on the value itself: parsing is parseValue's job, so an
// unparsable argument is told apart from an exhausted command line.
func nextArgValue(parser argvdeps.Parser) (string, bool) {
	raw, err := parser.GetNextStringArg()
	if err != nil {
		return "", false
	}
	return raw, true
}

// parseValue converts one raw command-line value to the type its declaration
// names, reporting the failure in the CLI's own words rather than the parser
// library's. A string is already what it will be, so only the numbers convert.
func parseValue(sandbox *api.Sandbox, subject string, name string, kind string, raw string) (any, bool) {
	switch kind {
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			sandbox.Deps.Std.Error("%s '%s': %q is not a valid integer\n", subject, name, raw)
			return nil, false
		}
		return value, true
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil {
			sandbox.Deps.Std.Error("%s '%s': %q is not a valid number\n", subject, name, raw)
			return nil, false
		}
		return value, true
	}
	return raw, true
}

// defaultValue is the declared default, spelled in entries.yaml as text, in
// the type the declaration names.
func defaultValue(sandbox *api.Sandbox, kind string, raw string) any {
	switch kind {
	case typeBoolean:
		return raw == "true"
	case typeInt:
		value, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			return 0
		}
		return value
	case typeFloat:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil {
			return float64(0)
		}
		return value
	}
	return raw
}

// inRange enforces the min/max bounds a numeric flag or arg declares. A value
// of any other type carries no bounds and passes untested.
func inRange(sandbox *api.Sandbox, subject string, name string, kind string, value any, min float64, hasMin bool, max float64, hasMax bool) bool {
	if !hasMin && !hasMax {
		return true
	}

	number, ok := numberOf(value)
	if !ok {
		return true
	}

	if hasMin && number < min {
		sandbox.Deps.Std.Error("%s '%s' must be >= %s\n", subject, name, numberLabel(sandbox, kind, min))
		return false
	}
	if hasMax && number > max {
		sandbox.Deps.Std.Error("%s '%s' must be <= %s\n", subject, name, numberLabel(sandbox, kind, max))
		return false
	}
	return true
}

// numberOf reads a bound value as the number a bound is compared against.
func numberOf(value any) (float64, bool) {
	if number, ok := value.(int); ok {
		return float64(number), true
	}
	if number, ok := value.(float64); ok {
		return number, true
	}
	return 0, false
}

// numberLabel spells a bound the way its declaration does, so the message
// quotes what entries.yaml says.
func numberLabel(sandbox *api.Sandbox, kind string, value float64) string {
	if kind == typeInt {
		return sandbox.Deps.Stringsdeps.FormatInt(int64(value), 10)
	}
	return sandbox.Deps.Stringsdeps.FormatFloat(value, 'g', -1, 64)
}

// checkUnknownFlags reports the first argument that still looks like a flag
// after every declared flag has been read — a typo such as --pathh, which
// would otherwise be ignored and leave the command running on a default.
func checkUnknownFlags(sandbox *api.Sandbox, parser argvdeps.Parser) bool {
	for i, used := range parser.Used {
		if used || !sandbox.Deps.Stringsdeps.HasPrefix(parser.Args[i], "-") {
			continue
		}
		sandbox.Deps.Std.Error("unknown flag %q — run '%s help' for the accepted flags\n", parser.Args[i], binaryName(sandbox))
		return false
	}
	return true
}

// checkUnusedArgs reports the first argument left over once every declared
// flag and positional arg has been read.
func checkUnusedArgs(sandbox *api.Sandbox, parser argvdeps.Parser) bool {
	for i, used := range parser.Used {
		if used {
			continue
		}
		sandbox.Deps.Std.Error("unexpected argument %q\n", parser.Args[i])
		return false
	}
	return true
}
