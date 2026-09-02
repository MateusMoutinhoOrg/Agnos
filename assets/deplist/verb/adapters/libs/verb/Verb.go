package verb

import (
	"{{.Module}}/sandbox/deps"
	argvdeps "{{.Module}}/sandbox/deps/argvdeps"
	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
)

// Bind fills deps.Deps.ArgvLib.New with the Verb argv-parser library's
// per-call constructor.
func Bind(deps *deps.Deps) {
	deps.Argvdeps.New = newArgvParser
}

// newArgvParser initializes the Verb argv-parser library over the given
// argument vector and copies it field by field onto the sandbox's local
// argvdeps.Parser.
func newArgvParser(args []string) argvdeps.Parser {
	inner := verblib.New(args)
	return argvdeps.Parser{
		Args: inner.Args,
		Used: inner.Used,

		IsPresent: inner.IsPresent,

		GetOptionsSize:   inner.GetOptionsSize,
		GetKeyValuesSize: inner.GetKeyValuesSize,

		GetStringOption:    inner.GetStringOption,
		GetIntOption:       inner.GetIntOption,
		GetDoubleOption:    inner.GetDoubleOption,
		GetTimestampOption: inner.GetTimestampOption,

		GetStringArg:    inner.GetStringArg,
		GetIntArg:       inner.GetIntArg,
		GetDoubleArg:    inner.GetDoubleArg,
		GetTimestampArg: inner.GetTimestampArg,

		GetNextStringArg:    inner.GetNextStringArg,
		GetNextIntArg:       inner.GetNextIntArg,
		GetNextDoubleArg:    inner.GetNextDoubleArg,
		GetNextTimestampArg: inner.GetNextTimestampArg,

		GetStringKeyValues:    inner.GetStringKeyValues,
		GetIntKeyValues:       inner.GetIntKeyValues,
		GetDoubleKeyValues:    inner.GetDoubleKeyValues,
		GetTimestampKeyValues: inner.GetTimestampKeyValues,
	}
}
