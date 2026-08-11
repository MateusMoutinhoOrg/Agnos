# `api.Lib.GetCategory`

**Type:** Field

## Signature

```go
GetCategory func(name string) (Category, bool)
```

## Description

Looks a stored category up by its unique name through the injected database and returns it as an [`api.Category`](/docs/LibUsage/References/PublicApi/api.Category.md), with the library's deps propagated in. The lookup goes through the database's key index, so it does not scan the collection.

Returns the **zero `Category`** and `false` when no category carries that name — always branch on the `bool`, since a struct return has no nil to compare against. Use [`AddCategory`](/docs/LibUsage/References/PublicApi/api.AddCategory.md) when you want the category created if it is missing.

The field holds a closure assigned by `GetCategoryFactory` over the `api.Lib` struct, so it reads `l.Deps` at call time.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `name` | `string` | The name the category was created with. |

## Returns

| Type | Description |
| :--- | :--- |
| [`api.Category`](/docs/LibUsage/References/PublicApi/api.Category.md) | The stored category, or the zero value on a miss. |
| `bool` | `false` when no category carries that name. |

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

	if groceries, ok := l.GetCategory("groceries"); ok {
		fmt.Println(groceries.String())
		return
	}
	fmt.Println("groceries: no such category")
}
```
