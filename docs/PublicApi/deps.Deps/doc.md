# `deps.Deps`

**Type:** Struct

## Definition

```go
type Deps struct {
	Argvdeps      argvdeps.Lib
	Embeddeps     embeddeps.Lib
	Goimportsdeps goimportsdeps.Lib
	Iodeps        iodeps.Lib
	Requestdeps   requestdeps.Lib
	Rundeps       rundeps.Lib
	Serializables serializables.Lib
	Std           std.Lib
}
```

## Description

The dependency contract every adapter must fill. A filled `Deps` is built by an adapter — see [`standard.New`](/docs/PublicApi/standard.New/doc.md) — and passed by pointer to [`sandbox.New`](/docs/PublicApi/sandbox.New/doc.md).

It carries **only** sub-contract structs, never a bare function field. Each one is the sandbox's copy of an outside library's api, declared under `sandbox/deps/<x>/` because the real thing — `os`, `net/http`, the Verb parser, `embed` — may not appear inside the sandbox; the adapter lib under `adapters/libs/` fills the copy. The loose runtime functions (`Now`, `Printf`, `Log`, `Error`, `Errorf`) are gathered into [`std.Lib`](/docs/PublicApi/std.Lib/doc.md) for the same reason. Sub-contracts whose real object is created per call expose that constructor as a field of their own `Lib` — [`requestdeps.Lib.NewRequest`](/docs/PublicApi/requestdeps.Lib/doc.md), [`argvdeps.Lib.New`](/docs/PublicApi/argvdeps.Lib/doc.md).

**Every field name is mechanical**: it is the title-cased sub-contract directory name (`sandbox/deps/iodeps/` → `Iodeps`), never a hand-picked alias, because `sandbox/deps/deps.go` is **generated** by `agnos build` from the listing of those directories. Installing a dep with `agnos dep-install` adds a directory and therefore a field; sandbox code must spell the field that way or stop compiling at the next build.

Because it is a struct and not an interface, a value returned by an adapter can be patched field by field before injection. The trade-off: the compiler cannot detect a field you forgot to fill — it stays nil and panics on first call. See [StructContracts](/docs/StructContracts/doc.md) and [ComposeDeps](/docs/ComposeDeps/doc.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Argvdeps` [`argvdeps.Lib`](/docs/PublicApi/argvdeps.Lib/doc.md) | A per-call argument-vector parser. How the generated CLI reads its command line. |
| `Embeddeps` [`embeddeps.Lib`](/docs/PublicApi/embeddeps.Lib/doc.md) | The assets compiled into the binary, and Go template rendering over them. How every generated file is produced. |
| `Goimportsdeps` [`goimportsdeps.Lib`](/docs/PublicApi/goimportsdeps.Lib/doc.md) | A Go source reader. Wired; not yet used by any action. |
| `Iodeps` [`iodeps.Lib`](/docs/PublicApi/iodeps.Lib/doc.md) | The filesystem. What SmartIO persists through. |
| `Requestdeps` [`requestdeps.Lib`](/docs/PublicApi/requestdeps.Lib/doc.md) | A per-call HTTP request. Wired; not yet used by any action. |
| `Rundeps` [`rundeps.Lib`](/docs/PublicApi/rundeps.Lib/doc.md) | Running an external program. How `build` hands a project to the Go toolchain. |
| `Serializables` [`serializables.Lib`](/docs/PublicApi/serializables.Lib/doc.md) | Generic JSON/YAML values. How every parsable config is read and written. |
| `Std` [`std.Lib`](/docs/PublicApi/std.Lib/doc.md) | The clock and the three output channels. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
)

func main() {
	d := agnosadapter.New()

	// Replace one sub-contract field before injection: record every program
	// the sandbox would run, without running it.
	var ran []string
	d.Rundeps.Run = func(props rundeps.RunProps) (rundeps.Result, error) {
		ran = append(ran, props.Program+" "+fmt.Sprint(props.Args))
		return rundeps.Result{ExitCode: 0}, nil
	}

	lib := agnoslib.New(&d)
	lib.Cli.CliMain([]string{"verify", "--path", "./my-tool"})
	fmt.Println(ran) // [go [mod tidy] go [build ./cmd/... ./sandbox/... ./adapters/...]]
}
```
