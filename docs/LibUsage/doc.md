# LibUsage

```bash
go get github.com/MateusMoutinhoOrg/Agnos@latest
```

```go
package main 
import (
	"github.com/MateusMoutinhoOrg/Agnos/adapters/availables/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)
func main(){
deps := standard.New()          // every adapter lib bound
lib := sandbox.New(&deps)       // *api.Sandbox: Actions + Cli

}

```

## Custom Deps

Patch fields **before** `sandbox.New(&deps)`; the binders capture the pointer. Start from `standard.New()`: an unfilled field is a nil func that panics on first call.

```go
deps := standard.New()
var out bytes.Buffer
deps.Std.Printf = func(f string, a ...any) (int, error) { return fmt.Fprintf(&out, f, a...) }
lib := sandbox.New(&deps)
```

A permanent mix is its own `adapters/availables/<name>/new.go` binding only the libs you want (`iodeps.Bind(&deps)`, ...). `standard/new.go` is regenerated on every build; other dirs under `availables/` are left alone. Wire it from a `cmd/<name>/main.go` of your own, since `cmd/main/main.go` is generated.
