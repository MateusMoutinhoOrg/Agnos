# Embedded Assets

## Description
How the library reads text, and later images and item templates, out of files instead of holding them in Go string constants: the assets under [`/assets/`](/docs/Development/References/Structure.md#assets), the [`embeddeps.Lib`](/docs/LibUsage/References/PublicApi/embeddeps.Lib.md) contract they arrive through, and what a consumer can do with them.

---

## Why Assets Are a Dependency

Reading a file is an OS-bound effect, and compiling one into a binary needs Go's `//go:embed` directive — a package-level, filesystem-bound mechanic. The sandbox may do neither, so assets follow the same route every other effect follows: a field on the [`Deps`](/docs/LibUsage/References/PublicApi/deps.Deps.md) contract, filled by an adapter outside the sandbox. See [SandboxIsolation.md](/docs/Development/References/SandboxIsolation.md).

```go
type Deps struct {
	Now       func() time.Time
	Printf    func(format string, a ...any) (n int, err error)
	VerbLib   verbdeps.Lib
	KeepLib   keepdeps.Lib
	EmbedDeps embeddeps.Lib // ← the assets
}
```

`EmbedDeps` is a whole library injected as one struct field, like `VerbLib` and `KeepLib` before it: three read-only functions, no getter, no bridging type. The sandbox asks for an asset by path; where those bytes come from is entirely the adapter's decision.

---

## What Lives in the Assets

Every piece of standing text the command-line interface displays, and nothing else:

| Asset | Description |
|-------|-------------|
| `version.txt` | The interface version reported by `agnos version` and `--version` |
| `cli/usage.txt` | The help screen, printed for `help`, for `--help`, and after any refused command line |
| `cli/messages/<name>.txt` | One file per line the interface can print, named after what it reports |

A message file is a `Printf` format, so the ones naming a value carry the verb — and the quotes — the value is rendered in:

```text
cli/messages/unknown-command.txt   →  unknown command "%s"
cli/messages/category-not-found.txt →  no category named "%s"
```

The result is that `sandbox/internal/cli/` holds no display text at all. Rewording the interface, or translating it, is editing files under `/assets/`; the Go code addresses them by path and never changes. Adding or editing one is [HandleAssets.md](/docs/Development/Tutorials/HandleAssets.md).

---

## Reading Assets as a Consumer

The adapter roots the asset tree at the directory passed to [`standard.New`](/docs/LibUsage/References/PublicApi/standard.New.md) — `"."` for everything — and every path is resolved against that root:

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
)

func main() {
	d := agnosadapter.New("trackerdata", ".")

	// One asset, whole.
	version, err := d.EmbedDeps.ReadFile("version.txt")
	if err != nil {
		fmt.Println("missing asset:", err)
		return
	}
	fmt.Print(string(version)) // v0.1.0

	// Every message the interface can print, discovered rather than listed.
	messages, _ := d.EmbedDeps.ListFiles("cli/messages")
	fmt.Println(len(messages), "messages") // 18 messages
}
```

Rooting the adapter at a subdirectory shifts every path with it, so a program interested only in the interface's text can address `usage.txt` directly:

```go
d := agnosadapter.New("trackerdata", "cli")
usage, _ := d.EmbedDeps.ReadFile("usage.txt")
```

---

## Serving Assets from Somewhere Else

Because the contract is three plain function fields, an adapter can back them with anything — a directory on disk that operators may edit, an archive, a translation service — and the library will not notice. Patching the field on a `deps.Deps` an adapter returned is enough, as long as it happens **before** `lib.New`; see [HandleDependencies.md](/docs/Development/Tutorials/HandleDependencies.md#overwriting-a-single-behavior).

```go
d := agnosadapter.New("trackerdata", ".")

// Serve the Portuguese wording, falling back to what ships in the binary.
shipped := d.EmbedDeps.ReadFile
d.EmbedDeps.ReadFile = func(path string) ([]byte, error) {
	if translated, err := os.ReadFile(filepath.Join("pt-BR", path)); err == nil {
		return translated, nil
	}
	return shipped(path)
}

l := agnoslib.New(d)
os.Exit(l.Sandboxmain(os.Args[1:]))
```

An asset that cannot be read is reported by path — `agnos: missing asset cli/usage.txt` — rather than printed as an empty line, because a missing asset is a packaging mistake and not a user mistake.

---

## Who Needs Assets Filled

Only [`Sandboxmain`](/docs/LibUsage/References/PublicApi/api.Sandboxmain.md) reads them. A program calling `AddCategory`, `AddSpend`, `ListTransactions` or `Balance` directly never touches an asset, so a hand-built `deps.Deps` may leave `EmbedDeps` zero — with the usual cost of a struct contract: calling into the interface afterwards panics on the nil `ReadFile` rather than failing to compile. See [StructContracts.md](/docs/Development/References/StructContracts.md#what-it-costs).
