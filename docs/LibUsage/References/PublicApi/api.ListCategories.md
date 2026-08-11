# `api.Lib.ListCategories`

**Type:** Field

## Signature

```go
ListCategories func() []Category
```

## Description

Reads every stored category back out of the injected database and returns them as [`api.Category`](/docs/LibUsage/References/PublicApi/api.Category.md) values, oldest first, each with the library's deps propagated in.

Returns an empty slice when the tracker holds no category, and `nil` when the injected database could not be read — callers that only range over the result need not tell the two apart.

The field holds a closure assigned by `ListCategoriesFactory` over the `api.Lib` struct, so every call re-reads the database rather than replaying a snapshot.

## Parameters

_None._

## Returns

| Type | Description |
| :--- | :--- |
| `[]api.Category` | Every stored category, oldest first. |

## Examples

```go
package main

import (
	"fmt"

	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	l := agnoslib.New(agnosadapter.New("trackerdata"))

	for _, category := range l.ListCategories() {
		fmt.Println(category.String())
	}
}
```
