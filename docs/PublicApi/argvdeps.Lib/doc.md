# `argvdeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
	New func(args []string) Parser
}

type Parser struct {
	Args []string
	Used []bool

	IsPresent func(flags []string) bool

	GetOptionsSize   func(flags []string) int
	GetKeyValuesSize func(prefixes []string) int

	GetStringOption    func(flags []string, occurrence int) (string, error)
	GetIntOption       func(flags []string, occurrence int) (int, error)
	GetDoubleOption    func(flags []string, occurrence int) (float64, error)
	GetTimestampOption func(flags []string, occurrence int) (time.Time, error)

	GetStringArg    func(index int) (string, error)
	GetIntArg       func(index int) (int, error)
	GetDoubleArg    func(index int) (float64, error)
	GetTimestampArg func(index int) (time.Time, error)

	GetNextStringArg    func() (string, error)
	GetNextIntArg       func() (int, error)
	GetNextDoubleArg    func() (float64, error)
	GetNextTimestampArg func() (time.Time, error)

	GetStringKeyValues    func(prefixes []string, occurrence int) (string, error)
	GetIntKeyValues       func(prefixes []string, occurrence int) (int, error)
	GetDoubleKeyValues    func(prefixes []string, occurrence int) (float64, error)
	GetTimestampKeyValues func(prefixes []string, occurrence int) (time.Time, error)
}
```

## Description

The argv-parser constructor injected whole as `Deps.Argvdeps`. A parser is bound to one argument vector, so it is created per call — `New(args)` — rather than injected once; what the sandbox holds is the one-field `Lib`. The package is the sandbox's **copy** of the [Verb](https://github.com/MateusMoutinhoOrg/Verb) library's api, restated field for field because the sandbox may not import a third-party module; the `verb` adapter lib initializes Verb over the given args and assigns its fields onto the copy. Installed by the `argvdeps` dep, which `cli-init` pulls in.

Every argument starts out unread; calling any `Get*` field or `IsPresent` marks what it matched as used, so whatever is left in `Args` is exactly the positionals nothing asked for. The two `*Size` fields count without marking. The generated dispatch reads flags first, then drains positionals with `GetNextStringArg`, then reports any unread `-`-prefixed argument as an unknown flag and any other leftover as unexpected — see [CommandDispatch](/docs/CommandDispatch/doc.md#the-dispatch-function).

## Fields

| Field | Description |
| :--- | :--- |
| `New func(args []string) Parser` | Builds a parser bound to `args`. |
| `Parser.Args`, `Parser.Used` | The vector and, index for index, what has been matched. Read-only. |
| `Parser.IsPresent(flags)` | Whether any spelling occurs unread; marks it used. Never fails. |
| `Parser.GetOptionsSize(flags)`, `GetKeyValuesSize(prefixes)` | Count matches without marking. Pair with the occurrence-indexed getters to iterate. |
| `Parser.Get{String,Int,Double,Timestamp}Option(flags, occurrence)` | The value after the occurrence-th match, marking both used. Errors when out of range or the flag has no value. |
| `Parser.Get{…}Arg(index)` | The argument at an absolute index, marked used. |
| `Parser.GetNext{…}Arg()` | The first still-unused argument, in order. Errors when exhausted. |
| `Parser.Get{…}KeyValues(prefixes, occurrence)` | The text after a `key=` prefix on the occurrence-th matching argument. |

## Examples

```go
verb := deps.Argvdeps.New([]string{"greet", "-n", "bob", "2"})

action, _ := verb.GetNextStringArg()                          // "greet"
name, _ := verb.GetStringOption([]string{"--name", "-n"}, 0)  // "bob"
times, _ := verb.GetNextIntArg()                              // 2
_, err := verb.GetNextStringArg()                             // err: everything is used
```
