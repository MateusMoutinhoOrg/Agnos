package verbdeps

import "time"

// Factory is the argv-parser constructor injected whole as the Deps.VerbLib
// field — the same mechanic as iodeps.Lib. A parser is created per call
// rather than injected once, so what the sandbox holds is this one-field
// struct; the adapter, which lives outside the sandbox, fills New over the
// embedded Verb library.
type Factory struct {
	// New builds an argv parser bound to the given arguments.
	New func(args []string) Lib
}

// This package is this library's *copy* of the embedded Verb argv-parser
// library's public api. The sandbox may not import the embedded library —
// that would be a third-party import — so it restates the shape it needs
// here, field for field. The adapter, which lives outside the sandbox, is
// what fills these structs from the real library.
//
// Copying is cheap precisely because the embedded library exposes structs
// of function fields instead of interfaces: an adapter assigns the real
// library's fields straight into the copy.

// Lib mirrors the embedded Verb library's api.Lib — an argument-vector
// (argv) parser. Every argument starts out unread; calling any Get* field
// or IsPresent marks the argument(s) it matched as used, so whatever is
// left over in Args is exactly the positional arguments nothing asked for.
// The two *Size fields are the exception: they count matches without ever
// marking anything used.
//
// Each getter family (Option, Arg, NextArg, KeyValues) is exposed once per
// supported value type: String (raw text), Int (base-10), Double (float64),
// and Timestamp (RFC 3339). A typed getter marks its match as used even
// when parsing then fails.
type Lib struct {
	// Args is the argument vector being parsed. Every index-based field
	// refers to positions in this slice. Treat it as read-only: mutating it
	// leaves Used out of sync.
	Args []string
	// Used tracks, index for index against Args, which arguments have
	// already been matched by a previous call. Treat it as read-only.
	Used []bool

	// IsPresent reports whether any of the given flag spellings (e.g.
	// []string{"-q", "--quiet"}) occurs in the unread portion of Args,
	// marking the matched argument used. It never fails: "not present" is a
	// valid outcome.
	IsPresent func(flags []string) bool

	// GetOptionsSize counts how many arguments equal one of the given flag
	// spellings, regardless of Used, and never mutates Used. Pair it with
	// GetStringOption to iterate occurrences 0..size-1.
	GetOptionsSize func(flags []string) int
	// GetKeyValuesSize counts how many arguments start with one of the given
	// key=value prefixes (the separator is part of the prefix), regardless
	// of Used, and never mutates Used.
	GetKeyValuesSize func(prefixes []string) int

	// GetStringOption returns the argument following the occurrence-th
	// (0-based) match of the given flag spellings, marking both as used. It
	// errors when occurrence is out of range or the flag has no value after
	// it.
	GetStringOption func(flags []string, occurrence int) (string, error)
	// GetIntOption behaves like GetStringOption, additionally parsing the
	// value as a base-10 integer.
	GetIntOption func(flags []string, occurrence int) (int, error)
	// GetDoubleOption behaves like GetStringOption, additionally parsing the
	// value as a 64-bit floating-point number.
	GetDoubleOption func(flags []string, occurrence int) (float64, error)
	// GetTimestampOption behaves like GetStringOption, additionally parsing
	// the value as an RFC 3339 timestamp.
	GetTimestampOption func(flags []string, occurrence int) (time.Time, error)

	// GetStringArg returns the argument at the given absolute index of Args
	// and marks it used. It errors when index is out of range.
	GetStringArg func(index int) (string, error)
	// GetIntArg behaves like GetStringArg, additionally parsing the argument
	// as a base-10 integer.
	GetIntArg func(index int) (int, error)
	// GetDoubleArg behaves like GetStringArg, additionally parsing the
	// argument as a 64-bit floating-point number.
	GetDoubleArg func(index int) (float64, error)
	// GetTimestampArg behaves like GetStringArg, additionally parsing the
	// argument as an RFC 3339 timestamp.
	GetTimestampArg func(index int) (time.Time, error)

	// GetNextStringArg returns the first still-unused argument, in order,
	// and marks it used — the leftover positional arguments, drained one
	// call at a time. It errors when every argument has been used.
	GetNextStringArg func() (string, error)
	// GetNextIntArg behaves like GetNextStringArg, additionally parsing the
	// argument as a base-10 integer.
	GetNextIntArg func() (int, error)
	// GetNextDoubleArg behaves like GetNextStringArg, additionally parsing
	// the argument as a 64-bit floating-point number.
	GetNextDoubleArg func() (float64, error)
	// GetNextTimestampArg behaves like GetNextStringArg, additionally
	// parsing the argument as an RFC 3339 timestamp.
	GetNextTimestampArg func() (time.Time, error)

	// GetStringKeyValues returns the text after the matched prefix of the
	// occurrence-th (0-based) argument starting with one of the given
	// key=value prefixes, marking it used. It errors when occurrence is out
	// of range or the value portion is empty.
	GetStringKeyValues func(prefixes []string, occurrence int) (string, error)
	// GetIntKeyValues behaves like GetStringKeyValues, additionally parsing
	// the value portion as a base-10 integer.
	GetIntKeyValues func(prefixes []string, occurrence int) (int, error)
	// GetDoubleKeyValues behaves like GetStringKeyValues, additionally
	// parsing the value portion as a 64-bit floating-point number.
	GetDoubleKeyValues func(prefixes []string, occurrence int) (float64, error)
	// GetTimestampKeyValues behaves like GetStringKeyValues, additionally
	// parsing the value portion as an RFC 3339 timestamp.
	GetTimestampKeyValues func(prefixes []string, occurrence int) (time.Time, error)
}
