package verb

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps/argvdeps"
	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
)

// NewArgvLib fills deps.Deps.ArgvLib.New: the Verb argv-parser library,
// initialized over the given argument vector and copied field by field onto
// the sandbox's local argvdeps.Parser.
func NewArgvLib(args []string) argvdeps.Parser {
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
