# Cache a Value

## Description
Covers storing a value in the cache with `Set` and reading it back with `Get`, including handling a cache miss. Installing and initializing the lib is covered by [LibInitialization.md](/docs/Tutorials/LibInitialization.md); inspecting an entry's expiry is covered by [InspectEntryExpiration.md](/docs/Tutorials/InspectEntryExpiration.md).

---

## Workflow
1. Initialize the lib as shown in [LibInitialization.md](/docs/Tutorials/LibInitialization.md):
   ```go
   deps := agnosadapter.New("cache.json")
   l := agnoslib.New(deps)
   ```
2. Store a value under a key with `Set`, passing the time-to-live in seconds:
   ```go
   // Keep "session-id" for 60 seconds
   l.Set("session-id", "abc123", 60)
   ```
3. Read the value back with `Get`. It returns an `api.Entry` and a boolean that is `false` when the key is absent **or** already expired — always check it before using the entry:
   ```go
   if entry, ok := l.Get("session-id"); ok {
       println(entry.Value)
   } else {
       println("cache miss: key absent or expired")
   }
   ```
4. To replace a value, call `Set` again with the same key — the new value and TTL overwrite the previous ones:
   ```go
   l.Set("session-id", "def456", 120)
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
    agnosadapter "github.com/MateusMoutinhoOrg/Agnos-Cli/adapters/standard"
    agnoslib "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox"
)

func main() {
    deps := agnosadapter.New("cache.json")
    l := agnoslib.New(deps)

    // Keep "session-id" for 60 seconds
    l.Set("session-id", "abc123", 60)

    if entry, ok := l.Get("session-id"); ok {
        println(entry.Value)
    } else {
        println("cache miss: key absent or expired")
    }

    // Replace the value — the new value and TTL overwrite the previous ones
    l.Set("session-id", "def456", 120)
}
```
