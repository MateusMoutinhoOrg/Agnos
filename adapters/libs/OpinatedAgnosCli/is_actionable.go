package opinatedagnoscli

import (
	"regexp"
	"strconv"
	"strings"

	opinatedagnoscli "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/OpinatedAgnosCli"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/argvdeps"
)

// endOfFlags is the token after which every token is a segment, whatever it
// starts with: `explain-command -- add-flag x --command c`.
const endOfFlags = "--"

// uuidPattern is what a `uuid` arg has to read as: the canonical 8-4-4-4-12
// hex form.
const uuidPattern = `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`

// isActionable reports whether one bound command — its Argv set — is for the
// command line it carries: the line has the Segments the command declares,
// every arg declaring a trigger finds its slice, converts to its Type and
// matches it, every arg that finds its slice converts, and every flag
// declaring a trigger brings a value that matches it. An arg or a flag that is
// merely missing is not a non-match: that is the binder's to answer, with a
// usage error.
//
// It is mirrored by agnos's own sandbox/internal/utils/command_match.go, which
// reads a command.yaml for `explain-command`: a change to one is a change to
// the other.
func isActionable(parsers argvdeps.Sandbox, command *opinatedagnoscli.Command) bool {
	segments, _ := splitArgv(command.Argv)
	if command.Segments > 0 && len(segments) != command.Segments {
		return false
	}

	for _, arg := range command.Args {
		values, found := argSlice(segments, arg)
		if !found {
			if arg.Trigger.Exist {
				return false
			}
			continue
		}
		if _, converts := argValue(arg, values); !converts {
			return false
		}
		if arg.Trigger.Exist && !matchTrigger(arg.Trigger, strings.Join(values, " "), true) {
			return false
		}
	}

	for _, flag := range command.Flags {
		if !flag.Trigger.Exist {
			continue
		}
		values := flagValues(parsers, command.Argv, flag)
		if len(values) == 0 || !matchTrigger(flag.Trigger, values[0], false) {
			return false
		}
	}

	return true
}

// splitArgv reads a command line into its segments and the index in argv of
// each: every token before the first flag — one starting with "-" that is not
// a number, so `sum -1 3` reads -1 as a segment — then every token after a bare
// "--", which is consumed with them. Everything between is the flags' to read.
func splitArgv(argv []string) ([]string, []int) {
	segments := []string{}
	indices := []int{}

	index := 0
	for ; index < len(argv); index++ {
		if isFlagToken(argv[index]) {
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

// isFlagToken reports whether a token of the command line is a flag: it starts
// with "-" and does not read as a number.
func isFlagToken(token string) bool {
	if !strings.HasPrefix(token, "-") || token == "-" {
		return false
	}
	_, err := strconv.ParseFloat(token, 64)
	return err != nil
}

// assignedKeys is the --key= prefix of every key of a flag, what the
// --key=value spelling of it starts with.
func assignedKeys(keys []string) []string {
	assigned := make([]string, 0, len(keys))
	for _, key := range keys {
		assigned = append(assigned, key+"=")
	}
	return assigned
}

// flagsEnd is the index of the bare "--" in argv, len(argv) when there is
// none: the flags are read before it and nowhere after.
func flagsEnd(argv []string) int {
	for index, token := range argv {
		if token == endOfFlags {
			return index
		}
	}
	return len(argv)
}

// argSlice is the segments one arg reads, Start to End with End -1 standing
// for the last one, and whether the command line has them. An arg reading to
// the last segment finds an empty slice on a line that stops right at Start.
func argSlice(segments []string, arg opinatedagnoscli.CommandArg) ([]string, bool) {
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

// argValue converts the slice one arg read to its Type: one segment to a
// string, an int, a float64 or a uuid string, several to a []string. It
// reports false for a slice that will not convert — a non-match.
func argValue(arg opinatedagnoscli.CommandArg, values []string) (any, bool) {
	if arg.End != arg.Start {
		return values, true
	}
	if len(values) != 1 {
		return nil, false
	}
	return convertArg(arg.Type, values[0])
}

// convertArg converts one segment to an arg type.
func convertArg(kind opinatedagnoscli.ArgType, text string) (any, bool) {
	switch kind {
	case opinatedagnoscli.IntegerArg:
		value, err := strconv.Atoi(text)
		return value, err == nil
	case opinatedagnoscli.NumberArg:
		value, err := strconv.ParseFloat(text, 64)
		return value, err == nil && isFinite(value)
	case opinatedagnoscli.UuidArg:
		matched, err := regexp.MatchString(uuidPattern, text)
		return text, err == nil && matched
	}
	return text, true
}

// isFinite reports whether a parsed number is one: NaN and ±Inf parse, and
// are not. Subtracting a number from itself gives 0 for every finite one and
// NaN for both.
func isFinite(value float64) bool {
	return value-value == 0
}

// flagValues is every raw value one flag brings, in order, read off a parser
// of its own so nothing is consumed: "true" once for a boolean flag that is
// present, one value per occurrence for any other. It is empty when the flag is
// absent.
func flagValues(parsers argvdeps.Sandbox, argv []string, flag opinatedagnoscli.CommandFlag) []string {
	parser := parsers.New(argv[:flagsEnd(argv)])
	if flag.Type == opinatedagnoscli.BooleanFlag {
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
	assigned := assignedKeys(flag.Keys)
	for occurrence := 0; occurrence < parser.GetKeyValuesSize(assigned); occurrence++ {
		value, err := parser.GetStringKeyValues(assigned, occurrence)
		if err != nil {
			break
		}
		values = append(values, value)
	}
	return values
}
