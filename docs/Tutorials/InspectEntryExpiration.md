# Inspect an Entry's Expiration

## Description
Covers reading the expiry information of a cached `api.Entry` with `ExpiresAt` and `IsExpired`. Storing and reading values is covered by [CacheValue.md](/docs/Tutorials/CacheValue.md); how the entry consults the injected clock is explained in [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

---

## Workflow
1. Store a value and get its entry back, as shown in [CacheValue.md](/docs/Tutorials/CacheValue.md):
   ```go
   l.Set("token", "abc123", 30)
   entry, ok := l.Get("token")
   if !ok {
       return
   }
   ```
2. Read the exact moment the entry stops being valid with `ExpiresAt`:
   ```go
   println("valid until:", entry.ExpiresAt.String())
   ```
3. Check whether the entry has already expired with `IsExpired`. The entry carries the injected clock, so the check is live — an entry held in a variable can expire between calls:
   ```go
   if entry.IsExpired() {
       println("entry expired, fetch a fresh value")
   }
   ```
4. To act on the remaining lifetime, compare `ExpiresAt` with the current time:
   ```go
   remaining := time.Until(entry.ExpiresAt)
   if remaining < 5*time.Second {
       // about to expire — renew it
       l.Set("token", "def456", 30)
   }
   ```
5. Run the program:
   ```bash
   go run main.go
   ```

---

## Full Code
```go
package main

import (
    "time"

    agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
    deps := agnosadapter.New("cache.json")
    l := agnoslib.New(deps)

    l.Set("token", "abc123", 30)
    entry, ok := l.Get("token")
    if !ok {
        return
    }

    println("valid until:", entry.ExpiresAt.String())

    if entry.IsExpired() {
        println("entry expired, fetch a fresh value")
    }

    remaining := time.Until(entry.ExpiresAt)
    if remaining < 5*time.Second {
        // about to expire — renew it
        l.Set("token", "def456", 30)
    }
}
```
