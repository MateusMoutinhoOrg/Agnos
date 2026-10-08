package opinatedagnoscli

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	opinatedagnoscli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinatedAgnosCli"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
)

// entriesArgument is the position of the Entries pointer among the parameters
// of an InternalPureHandler: func(props, entries, response) error.
const entriesArgument = 1

// entriesTag is the struct tag an Entries field names what it is bound to by.
const entriesTag = "id"

// fullCommandId is the id every Entries carries the whole command line under.
const fullCommandId = "FullCommand"

// handle binds the command line of one bound command onto a fresh Entries and
// runs the command's InternalPureHandler with it. The Entries type is the
// command package's own, so it is built, filled and called by reflection:
// every field is filled by its `id` tag — FullCommand, one per arg, one per
// flag.
//
// The order is the order command.yaml reads in: the args, then the flags. A
// value that will not bind is raised and the handler never runs. What the
// command read is marked on the chain's shared Consumed, and a strict command
// then refuses any token no command of the chain read. The handler is handed
// the command line's shared CommandProps first: what a middleware earlier in
// the chain set there is what it reads. A failure it returns — built by Fail —
// is raised on this command; any other error is returned as it is.
func (run *line) handle(command *opinatedagnoscli.Command) error {
	entries := newIn(command.InternalPureHandler, entriesArgument)
	if entries == nil || numField(entries) < 0 {
		return fmt.Errorf("command %s: InternalPureHandler is not a func(props *commandprops.CommandProps, entries *Entries, response *api.CommandResponse) error", command.Name)
	}

	values := map[string]any{fullCommandId: command.Argv}

	segments, indices := splitArgv(command.Argv)
	for _, arg := range command.Args {
		value, ok, err := run.bindArg(command, segments, indices, arg)
		if !ok {
			return err
		}
		if value != nil {
			values[arg.Id] = value
		}
	}

	parser := run.props.Argvdeps.New(command.Argv[:flagsEnd(command.Argv)])
	for _, flag := range command.Flags {
		value, ok, err := run.bindFlag(command, parser, flag)
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
	if end := flagsEnd(command.Argv); end < len(command.Argv) {
		command.Consumed[end] = true
	}

	if command.Strict {
		if ok, err := run.checkConsumed(command); !ok {
			return err
		}
	}

	for index := 0; index < numField(entries); index++ {
		id := fieldTag(entries, index, entriesTag)
		value, has := values[id]
		if id == "" || !has {
			continue
		}
		if err := setField(entries, index, value); err != nil {
			return fmt.Errorf("command %s: Entries.%s: %s", command.Name,
				fieldName(entries, index), err.Error())
		}
	}

	out := call(command.InternalPureHandler, []any{command.Props, entries, command.Response})
	if len(out) == 1 && out[0] != nil {
		if failure, is := out[0].(*opinatedagnoscli.CommandFailure); is {
			return run.raiseFailure(command, failure)
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
func (run *line) bindArg(command *opinatedagnoscli.Command, segments []string, indices []int, arg opinatedagnoscli.CommandArg) (any, bool, error) {
	slice, found := argSlice(segments, arg)
	if found && len(slice) > 0 {
		value, _ := argValue(arg, slice)
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
		return nil, false, run.raise(command, opinatedagnoscli.BadUsageFailure, opinatedagnoscli.ExitUsage, arg.Id,
			fmt.Sprintf("required arg '%s' not provided", arg.Id), "")
	}
	if arg.HasDefault {
		if arg.End != arg.Start {
			return []string{arg.Default}, true, nil
		}
		value, _ := convertArg(arg.Type, arg.Default)
		return value, true, nil
	}
	return nil, true, nil
}

// bindFlag reads one flag off the command line, through the parser whose Used
// becomes the command's consumed tokens. A value that will not convert, or
// breaks its bounds, its enum or its pattern, and a missing required flag, are
// raised as usage errors.
func (run *line) bindFlag(command *opinatedagnoscli.Command, parser argvdeps.Parser, flag opinatedagnoscli.CommandFlag) (any, bool, error) {
	if flag.Type == opinatedagnoscli.BooleanFlag {
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
			return nil, false, run.raise(command, opinatedagnoscli.BadUsageFailure, opinatedagnoscli.ExitUsage, flag.Id,
				fmt.Sprintf("flag '%s': expected a value after it", flag.Keys[0]), "")
		}
		raws = append(raws, raw)
	}
	// The --key=value spelling of the same flag.
	assigned := assignedKeys(flag.Keys)
	for occurrence := 0; occurrence < parser.GetKeyValuesSize(assigned); occurrence++ {
		raw, err := parser.GetStringKeyValues(assigned, occurrence)
		if err != nil {
			return nil, false, run.raise(command, opinatedagnoscli.BadUsageFailure, opinatedagnoscli.ExitUsage, flag.Id,
				fmt.Sprintf("flag '%s': expected a value after the =", flag.Keys[0]), "")
		}
		raws = append(raws, raw)
	}

	if len(raws) == 0 {
		if flag.Required {
			return nil, false, run.raise(command, opinatedagnoscli.BadUsageFailure, opinatedagnoscli.ExitUsage, flag.Id,
				fmt.Sprintf("required flag '%s' not provided", flag.Keys[0]), "")
		}
		if !flag.HasDefault {
			return nil, true, nil
		}
		raws = []string{flag.Default}
	}

	array := flag.Type == opinatedagnoscli.StringArrayFlag || flag.Type == opinatedagnoscli.IntegerArrayFlag
	if !array {
		raws = raws[:1]
	}

	texts := []string{}
	ints := []int{}
	var scalar any
	for _, raw := range raws {
		value, message := flagValue(flag, raw)
		if message != "" {
			return nil, false, run.raise(command, opinatedagnoscli.BadUsageFailure, opinatedagnoscli.ExitUsage, flag.Id,
				fmt.Sprintf("flag '%s': %s", flag.Keys[0], message), "")
		}
		switch typed := value.(type) {
		case int:
			ints = append(ints, typed)
		case string:
			texts = append(texts, typed)
		}
		scalar = value
	}

	switch flag.Type {
	case opinatedagnoscli.StringArrayFlag:
		return texts, true, nil
	case opinatedagnoscli.IntegerArrayFlag:
		return ints, true, nil
	}
	return scalar, true, nil
}

// flagValue converts one raw flag value to its type and checks it against the
// flag's bounds, enum and pattern. The message is "" when it passed.
func flagValue(flag opinatedagnoscli.CommandFlag, raw string) (any, string) {
	for _, accepted := range flag.Enum {
		if accepted == raw {
			break
		}
		if accepted == flag.Enum[len(flag.Enum)-1] {
			return nil, fmt.Sprintf("%q is not one of %s", raw, strings.Join(flag.Enum, ", "))
		}
	}
	if flag.Pattern != "" {
		matched, err := regexp.MatchString(flag.Pattern, raw)
		if err != nil || !matched {
			return nil, fmt.Sprintf("%q does not match %s", raw, flag.Pattern)
		}
	}

	var number float64
	var value any = raw
	switch flag.Type {
	case opinatedagnoscli.IntegerFlag, opinatedagnoscli.IntegerArrayFlag:
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return nil, fmt.Sprintf("%q is not a valid integer", raw)
		}
		value, number = parsed, float64(parsed)
	case opinatedagnoscli.NumberFlag:
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil || !isFinite(parsed) {
			return nil, fmt.Sprintf("%q is not a valid number", raw)
		}
		value, number = parsed, parsed
	default:
		return value, ""
	}

	if flag.HasMin && number < flag.Min {
		return nil, fmt.Sprintf("must be >= %s", strconv.FormatFloat(flag.Min, 'g', -1, 64))
	}
	if flag.HasMax && number > flag.Max {
		return nil, fmt.Sprintf("must be <= %s", strconv.FormatFloat(flag.Max, 'g', -1, 64))
	}
	return value, ""
}

// checkConsumed refuses, before a strict command runs, the first token no
// command of the chain read: one starting with "-" is an unknown flag — a
// typo such as --pathh, which would otherwise leave the command running on a
// default — and any other an unexpected argument.
func (run *line) checkConsumed(command *opinatedagnoscli.Command) (bool, error) {
	for index, token := range command.Argv {
		if command.Consumed[index] {
			continue
		}
		if isFlagToken(token) {
			return false, run.raise(command, opinatedagnoscli.UnknownFlagFailure, opinatedagnoscli.ExitUsage, token,
				fmt.Sprintf("unknown flag %q", token), "")
		}
		return false, run.raise(command, opinatedagnoscli.UnexpectedArgFailure, opinatedagnoscli.ExitUsage, token,
			fmt.Sprintf("unexpected argument %q", token), "")
	}
	return true, nil
}
