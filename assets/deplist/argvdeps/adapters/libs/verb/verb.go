package verb

import (
	"time"

	"{{.Module}}/sandbox/deps"
	argvdeps "{{.Module}}/sandbox/deps/argvdeps"
	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
)

// Bind fills deps.Deps.Argvdeps.New with the Verb argv-parser library's
// per-call constructor.
func Bind(deps *deps.Deps) {
	deps.Argvdeps.New = newArgvParser
}

// newArgvParser initializes the Verb argv-parser library over the given
// argument vector and copies it field by field onto the sandbox's local
// argvdeps.Parser. The Timestamp getters are the one family that cannot be
// assigned straight across: the sandbox may not name a `time.Time`, so each
// is wrapped by unixNano.
func newArgvParser(args []string) argvdeps.Parser {
	inner := verblib.New(args)
	return argvdeps.Parser{
		Args: inner.Args,
		Used: inner.Used,

		IsPresent: inner.IsPresent,

		GetOptionsSize:   inner.GetOptionsSize,
		GetKeyValuesSize: inner.GetKeyValuesSize,

		GetStringOption: inner.GetStringOption,
		GetIntOption:    inner.GetIntOption,
		GetDoubleOption: inner.GetDoubleOption,
		GetTimestampOption: func(flags []string, occurrence int) (int64, error) {
			return unixNano(inner.GetTimestampOption(flags, occurrence))
		},

		GetStringArg: inner.GetStringArg,
		GetIntArg:    inner.GetIntArg,
		GetDoubleArg: inner.GetDoubleArg,
		GetTimestampArg: func(index int) (int64, error) {
			return unixNano(inner.GetTimestampArg(index))
		},

		GetNextStringArg: inner.GetNextStringArg,
		GetNextIntArg:    inner.GetNextIntArg,
		GetNextDoubleArg: inner.GetNextDoubleArg,
		GetNextTimestampArg: func() (int64, error) {
			return unixNano(inner.GetNextTimestampArg())
		},

		GetStringKeyValues: inner.GetStringKeyValues,
		GetIntKeyValues:    inner.GetIntKeyValues,
		GetDoubleKeyValues: inner.GetDoubleKeyValues,
		GetTimestampKeyValues: func(prefixes []string, occurrence int) (int64, error) {
			return unixNano(inner.GetTimestampKeyValues(prefixes, occurrence))
		},
	}
}

// unixNano converts one Verb timestamp result into the nanoseconds since
// the Unix epoch the sandbox's argvdeps.Parser reports. A failed parse keeps
// its error and yields 0.
func unixNano(parsed time.Time, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return parsed.UnixNano(), nil
}
