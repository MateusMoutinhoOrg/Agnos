# `deps.Deps`

**Type:** Struct

## Definition

```go
type Deps struct {
	Now       func() time.Time
	Printf    func(format string, a ...any) (n int, err error)
	VerbLib   verbdeps.Lib
	KeepLib   keepdeps.Lib
	EmbedDeps embeddeps.Lib
}
```

## Description

The dependency contract every adapter must fill. Each field is one injectable behavior the financial tracker needs: `Now` is the clock (injecting it fixes the timestamp a category or transaction is stamped with, which is what makes the tracker testable), `Printf` is the writer the command-line interface reports through (injecting it captures the whole interface's output), `KeepLib` is the schema database every category and transaction is persisted in (an adapter can back it with the filesystem or with any other backend), and `EmbedDeps` is where every word the interface displays comes from. A filled `Deps` is built by an adapter — see [`standard.New`](/docs/LibUsage/References/PublicApi/standard.New.md) — and passed to [`lib.New`](/docs/LibUsage/References/PublicApi/lib.New.md).

`Printf` is the library's only way of emitting text: `Sandboxmain` prints every result, every error, and its usage screen through it, so the sandbox never touches a stream itself and the interface can be run against a buffer as easily as against a terminal. What it prints, though, is not written in the sandbox: the usage screen, the version, and each message come from `EmbedDeps`, so changing the interface's wording is editing an asset under `/assets/` rather than editing Go. Only `Sandboxmain` reads assets — a program calling `AddCategory`, `AddSpend`, or `Balance` directly never touches one, and can leave `EmbedDeps` unfilled.

`VerbLib`, `KeepLib` and `EmbedDeps` are the exceptions to "every field is a function": the dependency is itself a library built with this pattern, so it arrives as one plain struct field — [`verbdeps.Lib`](/docs/LibUsage/References/PublicApi/verbdeps.Lib.md), [`keepdeps.Lib`](/docs/LibUsage/References/PublicApi/keepdeps.Lib.md), [`embeddeps.Lib`](/docs/LibUsage/References/PublicApi/embeddeps.Lib.md) — with no getter around it. The sandbox never imports the embedded Verb or Keep libraries, nor Go's `embed` machinery; it declares a copy of each api in `sandbox/contracts/deps/verbdeps/`, `keepdeps/` and `embeddeps/`, and the adapter, which lives outside the sandbox, initializes the real library and assigns its fields onto that copy.

Because it is a struct and not an interface, a value returned by an adapter can be patched field by field before injection, and a custom contract needs no type declaration at all. The trade-off: the compiler cannot detect a field you forgot to fill — it stays nil and panics on first call. See [StructContracts.md](/docs/Development/References/StructContracts.md) and [HandleDependencies.md](/docs/Development/Tutorials/HandleDependencies.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Now func() time.Time` | Returns the current time, used to stamp categories and transactions as they are created. |
| `Printf func(format string, a ...any) (n int, err error)` | Writes one formatted message to the interface's output — the only way the library emits text. |
| `VerbLib verbdeps.Lib` | The embedded Verb argv parser, already initialized by the adapter over the argument vector that adapter chose. |
| `KeepLib keepdeps.Lib` | The embedded Keep schema database every category and transaction is stored in, already wired by the adapter to the storage backend that adapter chose. |
| `EmbedDeps embeddeps.Lib` | The embedded assets every piece of standing text the interface displays is read from — the usage screen, the version, and each message — already rooted by the adapter at the asset directory that adapter chose. |

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
