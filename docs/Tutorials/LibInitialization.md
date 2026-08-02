# Library Initialization

## Description
Covers installing the library and initializing it with the standard adapter in a new program. Using the cache after initialization is covered by [CacheValue.md](/docs/Tutorials/CacheValue.md). For other ways to build the dependencies, see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).

---

## Workflow
1. Install the lib:
   ```bash
   go get github.com/MateusMoutinhoOrg/Agnos@v0.0.7
   ```
2. Create a file called `main.go` with the following code:
   ```go
   package main

   // 1. Import the standard adapter and the lib
   import (
       agnosadapter "github.com/MateusMoutinhoOrg/Agnos/adapters/standard"
       agnoslib "github.com/MateusMoutinhoOrg/Agnos/sandbox"
   )

   func main() {
       // 2. Create deps via an adapter (the "opinionated" part)
       deps := agnosadapter.New("cache.json")

       // 3. Inject deps into the pure library — a key/value cache with TTL
       l := agnoslib.New(deps)

       // 4. Use the library — it never knows which adapter is behind the scenes
       l.Set("greeting", "hello world", 60)
       if entry, ok := l.Get("greeting"); ok {
           println(entry.Value)
       }
   }
   ```
3. Run the code:
   ```bash
   go run main.go
   ```
