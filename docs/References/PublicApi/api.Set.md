# `api.Lib.Set`

**Type:** Field

## Signature

```go
Set func(key string, value string, ttlSeconds int)
```

## Description

Stores `value` under `key`, computing the expiry as `ttlSeconds` after the injected clock's current time (`l.Deps.Now()`), then persisting it through `l.Deps.Store`. Because the clock is injected, expiry is deterministic under a custom `Deps`.

The field holds a closure assigned by `SetFactory` over the `api.Lib` struct, so it reads `l.Deps` at call time — calling `l.Set(...)` is indistinguishable from calling a method.

## Parameters

| Parameter | Type | Description |
| :--- | :--- | :--- |
| `key` | `string` | The key to store the value under. |
| `value` | `string` | The value to cache. |
| `ttlSeconds` | `int` | How many seconds from now the entry stays valid. |

## Returns

_None._

## Examples

```go
package main

import (
	agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
	agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
	l := agnoslib.New(agnosadapter.New("cache.json"))

	l.Set("session:42", "active", 300) // valid for 5 minutes
}
```
