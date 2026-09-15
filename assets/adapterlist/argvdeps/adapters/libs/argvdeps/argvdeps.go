package argvdeps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"{{.Module}}/sandbox/deps"
	argvdeps "{{.Module}}/sandbox/deps/argvdeps"
)

// timestampLayout is the format every Timestamp getter parses values with.
const timestampLayout = time.RFC3339

// Bind fills deps.Deps.Argvdeps.New with the per-call argv parser below. A
// parser is bound to one argument vector, so what the contract holds is a
// constructor rather than an already-built parser.
func Bind(deps *deps.Deps) {
	deps.Argvdeps.New = newParser
}

// newParser builds one argv parser over args. The returned value carries the
// argument vector and a same-length Used slice, and every function field is a
// closure over that same value: the copy handed back shares both slices, so a
// getter marking an argument used is visible through the copy.
func newParser(args []string) argvdeps.Parser {
	parser := argvdeps.Parser{
		Args: args,
		Used: make([]bool, len(args)),
	}

	parser.IsPresent = func(flags []string) bool {
		for index, arg := range parser.Args {
			if parser.Used[index] {
				continue
			}
			if matchesFlag(arg, flags) {
				parser.Used[index] = true
				return true
			}
		}
		return false
	}

	parser.GetOptionsSize = func(flags []string) int {
		count := 0
		for _, arg := range parser.Args {
			if matchesFlag(arg, flags) {
				count++
			}
		}
		return count
	}
	parser.GetKeyValuesSize = func(prefixes []string) int {
		count := 0
		for _, arg := range parser.Args {
			if _, ok := matchPrefix(arg, prefixes); ok {
				count++
			}
		}
		return count
	}

	parser.GetStringOption = func(flags []string, occurrence int) (string, error) {
		return optionValue(&parser, flags, occurrence)
	}
	parser.GetIntOption = func(flags []string, occurrence int) (int, error) {
		return parseInt(optionValue(&parser, flags, occurrence))
	}
	parser.GetDoubleOption = func(flags []string, occurrence int) (float64, error) {
		return parseDouble(optionValue(&parser, flags, occurrence))
	}
	parser.GetTimestampOption = func(flags []string, occurrence int) (int64, error) {
		return parseTimestamp(optionValue(&parser, flags, occurrence))
	}

	parser.GetStringArg = func(index int) (string, error) {
		return argValue(&parser, index)
	}
	parser.GetIntArg = func(index int) (int, error) {
		return parseInt(argValue(&parser, index))
	}
	parser.GetDoubleArg = func(index int) (float64, error) {
		return parseDouble(argValue(&parser, index))
	}
	parser.GetTimestampArg = func(index int) (int64, error) {
		return parseTimestamp(argValue(&parser, index))
	}

	parser.GetNextStringArg = func() (string, error) {
		return nextArgValue(&parser)
	}
	parser.GetNextIntArg = func() (int, error) {
		return parseInt(nextArgValue(&parser))
	}
	parser.GetNextDoubleArg = func() (float64, error) {
		return parseDouble(nextArgValue(&parser))
	}
	parser.GetNextTimestampArg = func() (int64, error) {
		return parseTimestamp(nextArgValue(&parser))
	}

	parser.GetStringKeyValues = func(prefixes []string, occurrence int) (string, error) {
		return keyValue(&parser, prefixes, occurrence)
	}
	parser.GetIntKeyValues = func(prefixes []string, occurrence int) (int, error) {
		return parseInt(keyValue(&parser, prefixes, occurrence))
	}
	parser.GetDoubleKeyValues = func(prefixes []string, occurrence int) (float64, error) {
		return parseDouble(keyValue(&parser, prefixes, occurrence))
	}
	parser.GetTimestampKeyValues = func(prefixes []string, occurrence int) (int64, error) {
		return parseTimestamp(keyValue(&parser, prefixes, occurrence))
	}

	return parser
}

// matchesFlag reports whether arg equals one of the given flag spellings.
func matchesFlag(arg string, flags []string) bool {
	for _, flag := range flags {
		if arg == flag {
			return true
		}
	}
	return false
}

// matchPrefix reports whether arg starts with one of the given key=value
// prefixes, returning the matched prefix.
func matchPrefix(arg string, prefixes []string) (string, bool) {
	for _, prefix := range prefixes {
		if strings.HasPrefix(arg, prefix) {
			return prefix, true
		}
	}
	return "", false
}

// optionValue backs every Option getter: it locates the occurrence-th
// (0-based) argument matching one of flags, regardless of Used, then marks
// both it and the argument after it used and returns that argument.
func optionValue(parser *argvdeps.Parser, flags []string, occurrence int) (string, error) {
	count := 0
	for index, arg := range parser.Args {
		if !matchesFlag(arg, flags) {
			continue
		}
		if count == occurrence {
			if index+1 >= len(parser.Args) {
				return "", fmt.Errorf("argvdeps: option %v: flag %q has no following value", flags, arg)
			}
			parser.Used[index] = true
			parser.Used[index+1] = true
			return parser.Args[index+1], nil
		}
		count++
	}
	return "", fmt.Errorf("argvdeps: option %v: occurrence %d not found (only %d present)", flags, occurrence, count)
}

// argValue backs every Arg getter: it validates index against Args, marks it
// used and returns the argument at that absolute position.
func argValue(parser *argvdeps.Parser, index int) (string, error) {
	if index < 0 || index >= len(parser.Args) {
		return "", fmt.Errorf("argvdeps: arg index %d out of range (have %d arguments)", index, len(parser.Args))
	}
	parser.Used[index] = true
	return parser.Args[index], nil
}

// nextArgValue backs every NextArg getter: it finds the first not-yet-used
// argument in order, marks it used and returns it.
func nextArgValue(parser *argvdeps.Parser) (string, error) {
	for index, used := range parser.Used {
		if !used {
			parser.Used[index] = true
			return parser.Args[index], nil
		}
	}
	return "", fmt.Errorf("argvdeps: no unused arguments remaining")
}

// keyValue backs every KeyValues getter: it locates the occurrence-th
// (0-based) argument starting with one of prefixes, marks it used and returns
// the text after the matched prefix.
func keyValue(parser *argvdeps.Parser, prefixes []string, occurrence int) (string, error) {
	count := 0
	for index, arg := range parser.Args {
		prefix, ok := matchPrefix(arg, prefixes)
		if !ok {
			continue
		}
		if count == occurrence {
			parser.Used[index] = true
			value := arg[len(prefix):]
			if value == "" {
				return "", fmt.Errorf("argvdeps: key/value %v: occurrence %d has an empty value", prefixes, occurrence)
			}
			return value, nil
		}
		count++
	}
	return "", fmt.Errorf("argvdeps: key/value %v: occurrence %d not found (only %d present)", prefixes, occurrence, count)
}

// parseInt parses the text one of the getters above matched as a base-10
// integer. A getter that already failed is reported unchanged, so the typed
// getters read as one expression.
func parseInt(text string, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	value, parse_err := strconv.Atoi(text)
	if parse_err != nil {
		return 0, fmt.Errorf("argvdeps: %q is not a valid integer: %w", text, parse_err)
	}
	return value, nil
}

// parseDouble parses the matched text as a 64-bit floating-point number.
func parseDouble(text string, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	value, parse_err := strconv.ParseFloat(text, 64)
	if parse_err != nil {
		return 0, fmt.Errorf("argvdeps: %q is not a valid number: %w", text, parse_err)
	}
	return value, nil
}

// parseTimestamp parses the matched text as an RFC 3339 timestamp and reports
// it as nanoseconds since the Unix epoch: the contract may not name a
// time.Time, so the conversion happens here, outside the sandbox.
func parseTimestamp(text string, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	value, parse_err := time.Parse(timestampLayout, text)
	if parse_err != nil {
		return 0, fmt.Errorf("argvdeps: %q is not a valid RFC3339 timestamp: %w", text, parse_err)
	}
	return value.UnixNano(), nil
}
