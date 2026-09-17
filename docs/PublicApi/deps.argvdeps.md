# `deps.Argvdeps`

`sandbox/deps/argvdeps`

## `Sandbox`

Sandbox is the argv-parser constructor injected whole as the Deps.ArgvLib field — the same mechanic as requestdeps.Sandbox. A parser is bound to one argument vector, so it is created per call rather than injected once: what the sandbox holds is this one-field struct, and the adapter — which lives outside the sandbox — fills New over a concrete argv-parser library.

| Field | Type | Description |
| --- | --- | --- |
| `New` | `func(args []string) Parser` | New builds an argv parser bound to the given arguments. |

## `Parser`

Parser mirrors the concrete argv-parser library's api.Lib — an argument-vector (argv) parser. Every argument starts out unread; calling any Get* field or IsPresent marks the argument(s) it matched as used, so whatever is left over in Args is exactly the positional arguments nothing asked for. The two *Size fields are the exception: they count matches without ever marking anything used. Each getter family (Option, Arg, NextArg, KeyValues) is exposed once per supported value type: String (raw text), Int (base-10), Double (float64), and Timestamp (RFC 3339 text, reported as nanoseconds since the Unix epoch, UTC — the sandbox may not name a `time.Time`). A typed getter marks its match as used even when parsing then fails.

| Field | Type | Description |
| --- | --- | --- |
| `Args` | `[]string` | Args is the argument vector being parsed. Every index-based field refers to positions in this slice. Treat it as read-only: mutating it leaves Used out of sync. |
| `Used` | `[]bool` | Used tracks, index for index against Args, which arguments have already been matched by a previous call. Treat it as read-only. |
| `IsPresent` | `func(flags []string) bool` | IsPresent reports whether any of the given flag spellings (e.g. []string{"-q", "--quiet"}) occurs in the unread portion of Args, marking the matched argument used. It never fails: "not present" is a valid outcome. |
| `GetOptionsSize` | `func(flags []string) int` | GetOptionsSize counts how many arguments equal one of the given flag spellings, regardless of Used, and never mutates Used. Pair it with GetStringOption to iterate occurrences 0..size-1. |
| `GetKeyValuesSize` | `func(prefixes []string) int` | GetKeyValuesSize counts how many arguments start with one of the given key=value prefixes (the separator is part of the prefix), regardless of Used, and never mutates Used. |
| `GetStringOption` | `func(flags []string, occurrence int) (string, error)` | GetStringOption returns the argument following the occurrence-th (0-based) match of the given flag spellings, marking both as used. It errors when occurrence is out of range or the flag has no value after it. |
| `GetIntOption` | `func(flags []string, occurrence int) (int, error)` | GetIntOption behaves like GetStringOption, additionally parsing the value as a base-10 integer. |
| `GetDoubleOption` | `func(flags []string, occurrence int) (float64, error)` | GetDoubleOption behaves like GetStringOption, additionally parsing the value as a 64-bit floating-point number. |
| `GetTimestampOption` | `func(flags []string, occurrence int) (int64, error)` | GetTimestampOption behaves like GetStringOption, additionally parsing the value as an RFC 3339 timestamp and reporting it as nanoseconds since the Unix epoch, UTC. |
| `GetStringArg` | `func(index int) (string, error)` | GetStringArg returns the argument at the given absolute index of Args and marks it used. It errors when index is out of range. |
| `GetIntArg` | `func(index int) (int, error)` | GetIntArg behaves like GetStringArg, additionally parsing the argument as a base-10 integer. |
| `GetDoubleArg` | `func(index int) (float64, error)` | GetDoubleArg behaves like GetStringArg, additionally parsing the argument as a 64-bit floating-point number. |
| `GetTimestampArg` | `func(index int) (int64, error)` | GetTimestampArg behaves like GetStringArg, additionally parsing the argument as an RFC 3339 timestamp and reporting it as nanoseconds since the Unix epoch, UTC. |
| `GetNextStringArg` | `func() (string, error)` | GetNextStringArg returns the first still-unused argument, in order, and marks it used — the leftover positional arguments, drained one call at a time. It errors when every argument has been used. |
| `GetNextIntArg` | `func() (int, error)` | GetNextIntArg behaves like GetNextStringArg, additionally parsing the argument as a base-10 integer. |
| `GetNextDoubleArg` | `func() (float64, error)` | GetNextDoubleArg behaves like GetNextStringArg, additionally parsing the argument as a 64-bit floating-point number. |
| `GetNextTimestampArg` | `func() (int64, error)` | GetNextTimestampArg behaves like GetNextStringArg, additionally parsing the argument as an RFC 3339 timestamp and reporting it as nanoseconds since the Unix epoch, UTC. |
| `GetStringKeyValues` | `func(prefixes []string, occurrence int) (string, error)` | GetStringKeyValues returns the text after the matched prefix of the occurrence-th (0-based) argument starting with one of the given key=value prefixes, marking it used. It errors when occurrence is out of range or the value portion is empty. |
| `GetIntKeyValues` | `func(prefixes []string, occurrence int) (int, error)` | GetIntKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as a base-10 integer. |
| `GetDoubleKeyValues` | `func(prefixes []string, occurrence int) (float64, error)` | GetDoubleKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as a 64-bit floating-point number. |
| `GetTimestampKeyValues` | `func(prefixes []string, occurrence int) (int64, error)` | GetTimestampKeyValues behaves like GetStringKeyValues, additionally parsing the value portion as an RFC 3339 timestamp and reporting it as nanoseconds since the Unix epoch, UTC. |

[every contract](doc.md)
