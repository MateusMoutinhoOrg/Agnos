# `deps.Deps`

**Type:** Struct

## Definition

```go
type Deps struct {
	Now        func() time.Time
	Printf     func(format string, a ...any) (n int, err error)
	VerbLib    verbdeps.Lib
	KeepLib    keepdeps.Lib
	EmbedDeps  embeddeps.Lib
	IoLib      iodeps.Lib
	NewRequest func(url string) serverdeps.Request
}
```

## Description

The dependency contract every adapter must fill. A filled `Deps` is built by an adapter — see [`standard.New`](/docs/References/PublicApi/standard.New.md) — and passed to [`lib.New`](/docs/References/PublicApi/lib.New.md).

Three fields are what the financial tracker is built from: `Now` is the clock (injecting it fixes the timestamp a category or transaction is stamped with, which is what makes the tracker testable), `Printf` is the writer the command-line interface reports through (injecting it captures the whole interface's output), and `KeepLib` is the schema database every category and transaction is persisted in (an adapter can back it with the filesystem or with any other backend). `VerbLib` is how the interface reads its command line.

`EmbedDeps`, `IoLib` and `NewRequest` are **standing capabilities**: the tracker never calls them, and the contract carries them anyway because this repository is a template. A derived library reads embedded assets, touches the filesystem, or speaks HTTP without designing a contract for it first, and the standard adapter already fills all three. They are offered, not required.

`Printf` is the library's only way of emitting text: `Sandboxmain` prints every result, every error, and its usage screen through it, so the sandbox never touches a stream itself and the interface can be run against a buffer as easily as against a terminal. What it prints is written in [`sandbox/config`](/docs/References/Structure.md#sandboxconfig) as compile-time constants, so changing the interface's wording is editing that package and every reference stays under the compiler's eye.

`VerbLib`, `KeepLib`, `EmbedDeps` and `IoLib` are the exceptions to "every field is a function": the dependency is itself a library built with this pattern, so it arrives as one plain struct field — [`verbdeps.Lib`](/docs/References/PublicApi/verbdeps.Lib.md), [`keepdeps.Lib`](/docs/References/PublicApi/keepdeps.Lib.md), [`embeddeps.Lib`](/docs/References/PublicApi/embeddeps.Lib.md), [`iodeps.Lib`](/docs/References/PublicApi/iodeps.Lib.md) — with no getter around it. `NewRequest` stays a function because a [`serverdeps.Request`](/docs/References/PublicApi/serverdeps.Request.md) is created per call rather than injected once. The sandbox never imports the embedded Verb or Keep libraries, Go's `embed` machinery, `os`, or `net/http`; it declares a copy of each api under `sandbox/contracts/deps/`, and the adapter, which lives outside the sandbox, initializes the real library and assigns its fields onto that copy.

Because it is a struct and not an interface, a value returned by an adapter can be patched field by field before injection, and a custom contract needs no type declaration at all. The trade-off: the compiler cannot detect a field you forgot to fill — it stays nil and panics on first call. See [StructContracts.md](/docs/References/StructContracts.md) and [HandleDependencies.md](/docs/Tutorials/HandleDependencies.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Now func() time.Time` | Returns the current time, used to stamp categories and transactions as they are created. |
| `Printf func(format string, a ...any) (n int, err error)` | Writes one formatted message to the interface's output — the only way the library emits text. |
| `VerbLib verbdeps.Lib` | The embedded Verb argv parser, already initialized by the adapter over the argument vector that adapter chose. |
| `KeepLib keepdeps.Lib` | The embedded Keep schema database every category and transaction is stored in, already wired by the adapter to the storage backend that adapter chose. |
| `EmbedDeps embeddeps.Lib` | The embedded assets under [`/assets/`](/docs/References/Structure.md#assets) — templates, long-form text, images — already rooted by the adapter at the asset directory that adapter chose. A standing capability: the tracker never reads one. |
| `IoLib iodeps.Lib` | The filesystem: reading, writing, and listing paths on disk. A standing capability: the tracker persists everything through `KeepLib` instead. |
| `NewRequest func(url string) serverdeps.Request` | Opens an HTTP request bound to `url`, handed back for the caller to configure and send. A standing capability: nothing the tracker does leaves the machine. |

## Examples

```go
package main

import (
	"fmt"
	"time"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	// Start from an adapter — it is what fills KeepLib, the one dependency
	// the tracker cannot work without — then patch the single behavior this
	// program wants to control.
	d := agnosadapter.New("trackerdata")

	frozen := time.Unix(0, 0)
	d.Now = func() time.Time { return frozen }

	l := agnoslib.New(d)

	l.AddCategory("groceries")
	transaction, _ := l.AddSpend("groceries", "coffee beans", 1290)

	fmt.Println(transaction.OccurredAt.Equal(frozen)) // true
}
```
