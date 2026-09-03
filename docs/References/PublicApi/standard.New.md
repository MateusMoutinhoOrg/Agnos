# `standard.New`

**Type:** Function

## Signature

```go
func New() deps.Deps
```

## Description

Builds the default [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md): an empty struct over which every adapter lib's `Bind` is run — `embeddeps`, `goimportsdeps`, `iodeps`, `requestdeps`, `rundeps`, `serializables`, `std`, `verb` — each filling its one field with the real thing: the process streams and clock, the Verb argv parser, the assets compiled into the binary, `os`, `net/http`, `os/exec`, `go/parser`, `yaml.v3`. The file `adapters/availables/standard/new.go` is **generated** by `agnos build` from the directories of `adapters/libs/`, so it always binds every lib the project has. It returns the contract struct by value; take its address to hand it to [`sandbox.New`](/docs/References/PublicApi/sandbox.New.md), and patch any field before doing so — see [ComposeDeps.md](/docs/Tutorials/ComposeDeps.md).

## Parameters

None.

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | Every field filled by the standard implementation. |

## Examples

```go
package main

import (
	"bytes"
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
)

func main() {
	deps := agnosadapter.New()

	// Replace one behavior before injection: capture what the CLI prints.
	var out bytes.Buffer
	deps.Std.Printf = func(format string, a ...any) (int, error) {
		return fmt.Fprintf(&out, format, a...)
	}

	lib := agnoslib.New(&deps)
	lib.Cli.CliMain([]string{"version"})
	fmt.Print(out.String()) // Version:0.0.4
}
```
