# `standard.New`

**Type:** Function

## Signature

```go
func New(filePath string) deps.Deps
```

## Description

Creates a [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) backed by the standard adapter: a single JSON file configured via `filePath` for storage, the real wall clock for `Now`, and the embedded Verb argv-parser library — initialized over the process's own command line, `os.Args[1:]` — for [`VerbLib`](/docs/References/PublicApi/verbdeps.Lib.md). Values survive across runs of the process. Initializing Verb here, outside the sandbox, is what keeps the sandbox free of the embedded library. The factory returns the **contract struct**, never the concrete `StandardAdapter`, so consumers stay decoupled from the implementation — each field is filled by a factory whose closure reads the adapter instance, which is how the adapter's state travels with the deps. For all shipped adapters, see [Adapters.md](/docs/References/Adapters.md).

## Parameters

| Name | Type | Description |
| :--- | :--- | :--- |
| `filePath` | `string` | Path to the JSON file where records should be persisted. |

## Returns

| Type | Description |
| :--- | :--- |
| [`deps.Deps`](/docs/References/PublicApi/deps.Deps.md) | A fully filled dependency contract, ready to be passed to `lib.New`. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	d := agnosadapter.New("cache.json")
	l := agnoslib.New(d)

	l.Set("greeting", "hello world", 60)
	if entry, ok := l.Get("greeting"); ok {
		fmt.Println(entry.Value) // hello world — still there on the next run
	}
}
```
