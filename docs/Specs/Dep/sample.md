# Dep Sample: `argvdeps`

The shipped dep installing the `argvdeps` sub-contract, implemented by the `verb` adapter lib over the external Verb module.

```text
assets/deplist/argvdeps/
├── sandbox/deps/argvdeps/argvdeps.go     # the contract: Lib{New}, Parser{…}
└── adapters/libs/verb/Verb.go            # Bind(deps) — copies verblib.New(args) onto argvdeps.Parser
```

```yaml
# assets/depsversion.yaml — one line per dep pulling an external module
argvdeps: github.com/MateusMoutinhoOrg/Verb@v0.0.1
```

```go
// assets/deplist/argvdeps/adapters/libs/verb/Verb.go — the module path is a template variable
import (
	"{{.Module}}/sandbox/deps"
	argvdeps "{{.Module}}/sandbox/deps/argvdeps"
	verblib "github.com/MateusMoutinhoOrg/Verb/sandbox"
)
```

After `agnos dep-install argvdeps`, the project holds `sandbox/deps/argvdeps/argvdeps.go` and `adapters/libs/verb/Verb.go`, its `go.mod` requires Verb, and the follow-up build regenerates `Deps.Argvdeps` and `verb.Bind(&deps)`.
