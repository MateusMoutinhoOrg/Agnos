package command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/generated/trigger"
)

// endOfFlags is the token after which every token is a segment, whatever it
// starts with: `explain-command -- add-flag x --command c`.
const endOfFlags = "--"

// uuidPattern is what a `uuid` arg has to read as: the canonical 8-4-4-4-12
// hex form.
const uuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

// IsActionable reports whether one bound command — its Argv set — is for the
// command line it carries: the line has the Segments the command declares,
// every arg declaring a trigger finds its slice, converts to its Type and
// matches it, every arg that finds its slice converts, and every flag
// declaring a trigger brings a value that matches it. An arg or a flag that is
// merely missing is not a non-match: that is CommandHandler's to answer, with
// a usage error.
func IsActionable(sandbox *api.Sandbox, command *api.Command) bool {
	segments, _ := SplitArgv(sandbox, command.Argv)
	if command.Segments > 0 && len(segments) != command.Segments {
		return false
	}

	for _, arg := range command.Args {
		values, found := ArgSlice(sandbox, segments, arg)
		if !found {
			if arg.Trigger.Exist {
				return false
			}
			continue
		}
		if _, converts := ArgValue(sandbox, arg, values); !converts {
			return false
		}
		if arg.Trigger.Exist && !trigger.MatchTrigger(sandbox, arg.Trigger, sandbox.Deps.Stringsdeps.Join(values, " "), true) {
			return false
		}
	}

	for _, flag := range command.Flags {
		if !flag.Trigger.Exist {
			continue
		}
		values := FlagValues(sandbox, command.Argv, flag)
		if len(values) == 0 || !trigger.MatchTrigger(sandbox, flag.Trigger, values[0], false) {
			return false
		}
	}

	return true
}

// SplitArgv reads a command line into its segments and the index in argv of
// each: every token before the first one starting with "-", then every token
// after a bare "--", which is consumed with them. Everything between is the
// flags' to read.
func SplitArgv(sandbox *api.Sandbox, argv []string) ([]string, []int) {
	segments := []string{}
	indices := []int{}

	index := 0
	for ; index < len(argv); index++ {
		if sandbox.Deps.Stringsdeps.HasPrefix(argv[index], "-") {
			break
		}
		segments = append(segments, argv[index])
		indices = append(indices, index)
	}

	for ; index < len(argv); index++ {
		if argv[index] != endOfFlags {
			continue
		}
		for rest := index + 1; rest < len(argv); rest++ {
			segments = append(segments, argv[rest])
			indices = append(indices, rest)
		}
		break
	}

	return segments, indices
}

// FlagsEnd is the index of the bare "--" in argv, len(argv) when there is
// none: the flags are read before it and nowhere after.
func FlagsEnd(argv []string) int {
	for index, token := range argv {
		if token == endOfFlags {
			return index
		}
	}
	return len(argv)
}

// ArgSlice is the segments one arg reads, Start to End with End -1 standing
// for the last one, and whether the command line has them. An arg reading to
// the last segment finds an empty slice on a line that stops right at Start.
func ArgSlice(sandbox *api.Sandbox, segments []string, arg api.CommandArg) ([]string, bool) {
	end := arg.End
	if end < 0 {
		if arg.Start >= len(segments) {
			return []string{}, arg.Start == len(segments)
		}
		end = len(segments) - 1
	}
	if arg.Start < 0 || arg.Start > end || end >= len(segments) {
		return nil, false
	}
	return segments[arg.Start : end+1], true
}

// ArgValue converts the slice one arg read to its Type: one segment to a
// string, an int, a float64 or a uuid string, several to a []string. It
// reports false for a slice that will not convert — a non-match.
func ArgValue(sandbox *api.Sandbox, arg api.CommandArg, values []string) (any, bool) {
	if arg.End != arg.Start {
		return values, true
	}
	if len(values) != 1 {
		return nil, false
	}
	return convertArg(sandbox, arg.Type, values[0])
}

// convertArg converts one segment to an arg type.
func convertArg(sandbox *api.Sandbox, kind api.ArgType, text string) (any, bool) {
	switch kind {
	case api.IntegerArg:
		value, err := sandbox.Deps.Stringsdeps.Atoi(text)
		return value, err == nil
	case api.NumberArg:
		value, err := sandbox.Deps.Stringsdeps.ParseFloat(text, 64)
		return value, err == nil
	case api.UuidArg:
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(uuidPattern, text)
		return text, err == nil && matched
	}
	return text, true
}

// FlagValues is every raw value one flag brings, in order, read off a parser
// of its own so nothing is consumed: "true" once for a boolean flag that is
// present, one value per occurrence for any other. It is empty when the flag is
// absent.
func FlagValues(sandbox *api.Sandbox, argv []string, flag api.CommandFlag) []string {
	parser := sandbox.Deps.Argvdeps.New(argv[:FlagsEnd(argv)])
	if flag.Type == api.BooleanFlag {
		if parser.GetOptionsSize(flag.Keys) > 0 {
			return []string{"true"}
		}
		return []string{}
	}

	values := []string{}
	for occurrence := 0; occurrence < parser.GetOptionsSize(flag.Keys); occurrence++ {
		value, err := parser.GetStringOption(flag.Keys, occurrence)
		if err != nil {
			break
		}
		values = append(values, value)
	}
	return values
}
