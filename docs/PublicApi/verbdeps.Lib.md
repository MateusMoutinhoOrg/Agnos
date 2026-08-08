# `verbdeps.Lib`

**Type:** Struct

## Definition

```go
type Lib struct {
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

The sandbox's **copy** of the embedded Verb argv-parser library's `api.Lib`, declared in `sandbox/contracts/deps/verbdeps/` and injected whole as the [`deps.Deps.VerbLib`](/docs/PublicApi/deps.Deps.md) field. The sandbox may not import Verb — that would be a third-party import — so it restates the shape it needs, field for field; the adapter, which lives outside the sandbox, initializes the real library and assigns its fields onto this copy. Copying is cheap precisely because both sides are structs of function fields rather than interfaces: the copy is plain field assignment. See [StructContracts.md](/docs/StructContracts.md).

Every argument starts out unread. Calling any `Get*` field or `IsPresent` marks the argument(s) it matched as used in `Used`, so whatever is left in `Args` is exactly the positional arguments nothing asked for — drain them in order with `GetNextStringArg`. The two `*Size` fields are the exception: they count matches without ever marking anything used.

Each getter family (Option, Arg, NextArg, KeyValues) is exposed once per value type: `String` (raw text), `Int` (base-10), `Double` (`float64`), and `Timestamp` (RFC 3339, e.g. `"2024-01-02T15:04:05Z"`). A typed getter marks its match as used even when parsing then fails — the argument was still read, only its value turned out malformed.

## Fields

| Field | Description |
| :--- | :--- |
| `Args []string` | The argument vector being parsed. Read-only: mutating it leaves `Used` out of sync. |
| `Used []bool` | Index for index against `Args`, which arguments have already been matched. Read-only. |
| `IsPresent func(flags []string) bool` | Reports whether any of the given flag spellings occurs in the unread portion of `Args`, marking the match used. Never fails. |
| `GetOptionsSize func(flags []string) int` | Counts arguments equal to one of the flag spellings, regardless of `Used`, without marking anything. |
| `GetKeyValuesSize func(prefixes []string) int` | Counts arguments starting with one of the `key=` prefixes, regardless of `Used`, without marking anything. |
| `GetStringOption func(flags []string, occurrence int) (string, error)` | Returns the argument following the occurrence-th match of the flags, marking both used. |
| `GetIntOption` / `GetDoubleOption` / `GetTimestampOption` | As `GetStringOption`, additionally parsing the value as an int, a float64, or an RFC 3339 timestamp. |
| `GetStringArg func(index int) (string, error)` | Returns the argument at an absolute index of `Args` and marks it used. |
| `GetIntArg` / `GetDoubleArg` / `GetTimestampArg` | As `GetStringArg`, additionally parsing the argument. |
| `GetNextStringArg func() (string, error)` | Returns the first still-unused argument, in order, and marks it used. Errors once all are used. |
| `GetNextIntArg` / `GetNextDoubleArg` / `GetNextTimestampArg` | As `GetNextStringArg`, additionally parsing the argument. |
| `GetStringKeyValues func(prefixes []string, occurrence int) (string, error)` | Returns the value portion of the occurrence-th argument starting with one of the `key=` prefixes, marking it used. |
| `GetIntKeyValues` / `GetDoubleKeyValues` / `GetTimestampKeyValues` | As `GetStringKeyValues`, additionally parsing the value portion. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	// The adapter initializes the embedded Verb library and hands it back as
	// one field of the deps contract.
	d := agnosadapter.New("trackerdata")

	quiet := d.VerbLib.IsPresent([]string{"-q", "--quiet"})

	out, err := d.VerbLib.GetStringOption([]string{"-o", "--out"}, 0)
	if err != nil {
		fmt.Println("no output option:", err)
		return
	}

	count, err := d.VerbLib.GetIntKeyValues([]string{"count="}, 0)
	if err != nil {
		count = 1
	}

	// Everything never explicitly asked for is still unread.
	file, err := d.VerbLib.GetNextStringArg()
	if err != nil {
		fmt.Println("no file given")
		return
	}

	fmt.Println(out, count, file, quiet) // file.txt 7 report.md false
}
```
