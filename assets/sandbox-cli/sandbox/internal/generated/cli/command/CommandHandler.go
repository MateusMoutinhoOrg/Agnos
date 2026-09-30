package command

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/argvdeps"
	"{{.Module}}/sandbox/internal/generated/cliio"
)

// entriesArgument is the position of the Entries pointer among the parameters
// of an InternalPureHandler: func(props, entries, response) error.
const entriesArgument = 1

// entriesTag is the struct tag an Entries field names what it is bound to by.
const entriesTag = "id"

// fullCommandId is the id every Entries carries the whole command line under.
const fullCommandId = "FullCommand"

// CommandHandler binds the command line of one bound command onto a fresh
// Entries and runs the command's InternalPureHandler with it. The Entries type
// is the command package's own, so it is built, filled and called through
// Deps.Reflectdeps: every field is filled by its `id` tag — FullCommand, one
// per arg, one per flag.
//
// The order is the order command.yaml reads in: the args, then the flags. A
// value that will not bind is raised through cliio.Raise and the handler never
// runs. What the command read is marked on the chain's shared Consumed, and a
// strict command then refuses any token no command of the chain read. The
// handler is handed the command line's shared CommandProps first: what a
// middleware earlier in the chain set there is what it reads. A failure it
// returns — built by cliio.Fail — is raised on this command; any other error
// is returned as it is.
func CommandHandler(sandbox *api.Sandbox, command *api.Command) error {
	entries := sandbox.Deps.Reflectdeps.NewIn(command.InternalPureHandler, entriesArgument)
	if entries == nil || sandbox.Deps.Reflectdeps.NumField(entries) < 0 {
		return sandbox.Deps.Std.Errorf("command %s: InternalPureHandler is not a func(props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error", command.Name)
	}

	values := map[string]any{fullCommandId: command.Argv}

	segments, indices := SplitArgv(sandbox, command.Argv)
	for _, arg := range command.Args {
		value, ok, err := bindArg(sandbox, command, segments, indices, arg)
		if !ok {
			return err
		}
		if value != nil {
			values[arg.Id] = value
		}
	}

	parser := sandbox.Deps.Argvdeps.New(command.Argv[:FlagsEnd(command.Argv)])
	for _, flag := range command.Flags {
		value, ok, err := bindFlag(sandbox, command, parser, flag)
		if !ok {
			return err
		}
		if value != nil {
			values[flag.Id] = value
		}
	}
	for index, used := range parser.Used {
		if used {
			command.Consumed[index] = true
		}
	}
	if end := FlagsEnd(command.Argv); end < len(command.Argv) {
		command.Consumed[end] = true
	}

	if command.Strict {
		if ok, err := checkConsumed(sandbox, command); !ok {
			return err
		}
	}

	for index := 0; index < sandbox.Deps.Reflectdeps.NumField(entries); index++ {
		id := sandbox.Deps.Reflectdeps.FieldTag(entries, index, entriesTag)
		value, has := values[id]
		if id == "" || !has {
			continue
		}
		if err := sandbox.Deps.Reflectdeps.SetField(entries, index, value); err != nil {
			return sandbox.Deps.Std.Errorf("command %s: Entries.%s: %s", command.Name,
				sandbox.Deps.Reflectdeps.FieldName(entries, index), err.Error())
		}
	}

	out := sandbox.Deps.Reflectdeps.Call(command.InternalPureHandler, []any{command.Props, entries, command.Response})
	if len(out) == 1 && out[0] != nil {
		if failure, is := out[0].(*api.CommandFailure); is {
			return cliio.RaiseFailure(sandbox, command, failure)
		}
		if handler_error, is := out[0].(error); is {
			return handler_error
		}
	}
	return nil
}

// bindArg reads one arg off the segments and marks them consumed — on a strict
// command alone: a middleware reads the segments to match on them, and the
// command after it still has to declare every one it takes. A missing
// required arg is raised as a usage error; a missing optional one binds its
// default, or nothing. It reports false, with what the failure returned, when
// the handler must not run.
func bindArg(sandbox *api.Sandbox, command *api.Command, segments []string, indices []int, arg api.CommandArg) (any, bool, error) {
	slice, found := ArgSlice(sandbox, segments, arg)
	if found && len(slice) > 0 {
		value, _ := ArgValue(sandbox, arg, slice)
		if !command.Strict {
			return value, true, nil
		}
		end := arg.End
		if end < 0 {
			end = len(segments) - 1
		}
		for index := arg.Start; index <= end; index++ {
			command.Consumed[indices[index]] = true
		}
		return value, true, nil
	}

	if arg.Required {
		return nil, false, cliio.Raise(sandbox, command, api.BadUsageFailure, api.ExitUsage, arg.Id,
			sandbox.Deps.Std.Sprintf("required arg '%s' not provided", arg.Id))
	}
	if arg.HasDefault {
		if arg.End != arg.Start {
			return []string{arg.Default}, true, nil
		}
		value, _ := convertArg(sandbox, arg.Type, arg.Default)
		return value, true, nil
	}
	return nil, true, nil
}

// bindFlag reads one flag off the command line, through the parser whose Used
// becomes the command's consumed tokens. A value that will not convert, or
// breaks its bounds, its enum or its pattern, and a missing required flag, are
// raised as usage errors.
func bindFlag(sandbox *api.Sandbox, command *api.Command, parser argvdeps.Parser, flag api.CommandFlag) (any, bool, error) {
	if flag.Type == api.BooleanFlag {
		present := false
		for parser.IsPresent(flag.Keys) {
			present = true
		}
		return present, true, nil
	}

	raws := []string{}
	for occurrence := 0; occurrence < parser.GetOptionsSize(flag.Keys); occurrence++ {
		raw, err := parser.GetStringOption(flag.Keys, occurrence)
		if err != nil {
			return nil, false, cliio.Raise(sandbox, command, api.BadUsageFailure, api.ExitUsage, flag.Id,
				sandbox.Deps.Std.Sprintf("flag '%s': expected a value after it", flag.Keys[0]))
		}
		raws = append(raws, raw)
	}
	// The --key=value spelling of the same flag.
	assigned := AssignedKeys(flag.Keys)
	for occurrence := 0; occurrence < parser.GetKeyValuesSize(assigned); occurrence++ {
		raw, err := parser.GetStringKeyValues(assigned, occurrence)
		if err != nil {
			return nil, false, cliio.Raise(sandbox, command, api.BadUsageFailure, api.ExitUsage, flag.Id,
				sandbox.Deps.Std.Sprintf("flag '%s': expected a value after the =", flag.Keys[0]))
		}
		raws = append(raws, raw)
	}

	if len(raws) == 0 {
		if flag.Required {
			return nil, false, cliio.Raise(sandbox, command, api.BadUsageFailure, api.ExitUsage, flag.Id,
				sandbox.Deps.Std.Sprintf("required flag '%s' not provided", flag.Keys[0]))
		}
		if !flag.HasDefault {
			return nil, true, nil
		}
		raws = []string{flag.Default}
	}

	array := flag.Type == api.StringArrayFlag || flag.Type == api.IntegerArrayFlag
	if !array {
		raws = raws[:1]
	}

	strings := []string{}
	ints := []int{}
	var scalar any
	for _, raw := range raws {
		value, message := flagValue(sandbox, flag, raw)
		if message != "" {
			return nil, false, cliio.Raise(sandbox, command, api.BadUsageFailure, api.ExitUsage, flag.Id,
				sandbox.Deps.Std.Sprintf("flag '%s': %s", flag.Keys[0], message))
		}
		switch typed := value.(type) {
		case int:
			ints = append(ints, typed)
		case string:
			strings = append(strings, typed)
		}
		scalar = value
	}

	switch flag.Type {
	case api.StringArrayFlag:
		return strings, true, nil
	case api.IntegerArrayFlag:
		return ints, true, nil
	}
	return scalar, true, nil
}

// flagValue converts one raw flag value to its type and checks it against the
// flag's bounds, enum and pattern. The message is "" when it passed.
func flagValue(sandbox *api.Sandbox, flag api.CommandFlag, raw string) (any, string) {
	for _, accepted := range flag.Enum {
		if accepted == raw {
			break
		}
		if accepted == flag.Enum[len(flag.Enum)-1] {
			return nil, sandbox.Deps.Std.Sprintf("%q is not one of %s", raw, sandbox.Deps.Stringsdeps.Join(flag.Enum, ", "))
		}
	}
	if flag.Pattern != "" {
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(flag.Pattern, raw)
		if err != nil || !matched {
			return nil, sandbox.Deps.Std.Sprintf("%q does not match %s", raw, flag.Pattern)
		}
	}

	var number float64
	var value any = raw
	switch flag.Type {
	case api.IntegerFlag, api.IntegerArrayFlag:
		parsed, err := sandbox.Deps.Stringsdeps.Atoi(raw)
		if err != nil {
			return nil, sandbox.Deps.Std.Sprintf("%q is not a valid integer", raw)
		}
		value, number = parsed, float64(parsed)
	case api.NumberFlag:
		parsed, err := sandbox.Deps.Stringsdeps.ParseFloat(raw, 64)
		if err != nil || !IsFinite(parsed) {
			return nil, sandbox.Deps.Std.Sprintf("%q is not a valid number", raw)
		}
		value, number = parsed, parsed
	default:
		return value, ""
	}

	if flag.HasMin && number < flag.Min {
		return nil, sandbox.Deps.Std.Sprintf("must be >= %s", sandbox.Deps.Stringsdeps.FormatFloat(flag.Min, 'g', -1, 64))
	}
	if flag.HasMax && number > flag.Max {
		return nil, sandbox.Deps.Std.Sprintf("must be <= %s", sandbox.Deps.Stringsdeps.FormatFloat(flag.Max, 'g', -1, 64))
	}
	return value, ""
}

// checkConsumed refuses, before a strict command runs, the first token no
// command of the chain read: one starting with "-" is an unknown flag — a
// typo such as --pathh, which would otherwise leave the command running on a
// default — and any other an unexpected argument.
func checkConsumed(sandbox *api.Sandbox, command *api.Command) (bool, error) {
	for index, token := range command.Argv {
		if command.Consumed[index] {
			continue
		}
		if IsFlagToken(sandbox, token) {
			return false, cliio.Raise(sandbox, command, api.UnknownFlagFailure, api.ExitUsage, token,
				sandbox.Deps.Std.Sprintf("unknown flag %q", token))
		}
		return false, cliio.Raise(sandbox, command, api.UnexpectedArgFailure, api.ExitUsage, token,
			sandbox.Deps.Std.Sprintf("unexpected argument %q", token))
	}
	return true, nil
}
