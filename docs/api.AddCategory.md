# `api.Lib.AddCategory`

**Type:** Field

## Signature

```go
AddCategory func(name string) (Category, bool)
```

## Description

Creates the category `name` in the injected database, stamping it with the injected clock (`l.Deps.Now()`), and returns it as an [`api.Category`](/docs/api.Category.md) with the library's deps propagated in.

Creation is **idempotent**: the category name is a unique key, so a name already taken makes the injected database report a key conflict, and the stored category is returned instead of a failure. Re-running a program that starts by declaring its categories is therefore safe.

The field holds a closure assigned by `AddCategoryFactory` over the `api.Lib` struct, so it reads `l.Deps` at call time — calling `l.AddCategory(...)` is indistinguishable from calling a method.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | The category's unique name. Must not be empty. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Category`](/docs/api.Category.md) | The created or already-stored category, or the zero value on failure. |
| `bool` | `false` when `name` is empty or the record could not be written. |

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

	groceries, ok := l.AddCategory("groceries")
	if !ok {
		fmt.Println("could not create the category")
		return
	}
	fmt.Println(groceries.Name)
}
```
