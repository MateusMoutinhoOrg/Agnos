# DepList

The catalogue `{{.GeneratorName}} add-dep <dep>` installs from. A **dep** is the contract, one
directory under `sandbox/deps/<dep>/`, filling the `Deps` field of the same name title-cased.
An **adapter** is one implementation of it, one directory under `adapters/libs/<adapter>/`.
The two are separate catalogues because one contract may have several implementations — how
that works is [Adapters](../Adapters/doc.md). Signatures of the contracts already installed are
in [PublicApi](../PublicApi/doc.md#dependency-contracts).

| Dep | `Deps` field | Adapters | Backed by | Provides |
|---|---|---|---|---|
| `argvdeps` | `Argvdeps` | `verb` | `github.com/MateusMoutinhoOrg/Verb` (pinned in `go.mod`) | Per-call argv parser. Installed by `cli-init` |
| `dbdeps` | `Dbdeps` | `keep` | `github.com/MateusMoutinhoOrg/Keep` (pinned) | Schema database |
| `embeddeps` | `Embeddeps` | `embeddeps` + `assets/asset.go` | `embed`, `text/template` | Read and render files compiled into the binary |
| `goimportsdeps` | `Goimportsdeps` | `goimportsdeps` | `go/parser` | Go source reader (package, imports, declarations) |
| `hashdeps` | `Hashdeps` | `hashdeps` | `crypto/sha256`, `encoding/hex` | SHA-256 of a byte slice, lower-case hex |
| `iodeps` | `Iodeps` | `iodeps` | `os`, `path/filepath` | Filesystem. `WriteFile` creates parents; `RemoveDir` removes files too; `Join`/`Dir` build host paths |
| `requestdeps` | `Requestdeps` | `requestdeps` | `net/http` (30s timeout) | Per-call HTTP request |
| `rundeps` | `Rundeps` | `rundeps` | `os/exec` | Run a program to completion; stdout+stderr merged; non-zero exit is `Result.ExitCode`, not an error |
| `serializables` | `Serializables` | `serializables` | `gopkg.in/yaml.v3`, `encoding/json` | Generic JSON/YAML values |
| `serverdeps` | `Serverdeps` | `serverdeps` | `net/http` | Http server: opens the port, applies timeouts, hands every request to one handler. Installed by `server-init` |
| `sortdeps` | `Sortdeps` | `sortdeps`, `reflectsort` | `sort` / `reflect` | Sort a string slice, or any slice by a less function |
| `std` | `Std` | `std` | `time`, `fmt`, `runtime`, `os.Stdout/Stderr` | Clock, `Sprintf`, the host `Goos` and the three output channels. Installed by `cli-init` |
| `stringsdeps` | `Stringsdeps` | `stringsdeps` | `strings`, `strconv` | Text manipulation and string/number conversion |
| `templatedeps` | `Templatedeps` | `templatedeps` | `text/template` | Parse and execute one template over vars, with native funcs |

The first adapter of each row is the dep's `default-adapter`, the one `add-dep` installs when
`--adapter` names no other.

```bash
{{.GeneratorName}} list-deps                       # this table, plus what is installed here
{{.GeneratorName}} list-adapters                   # one row per adapter, and which available binds it
{{.GeneratorName}} add-dep sortdeps                # contract + default-adapter
{{.GeneratorName}} add-dep sortdeps --adapter reflectsort
{{.GeneratorName}} remove-dep sortdeps             # refused while an adapter fills it
{{.GeneratorName}} remove-dep sortdeps --with-adapters
```

Writing a contract of your own instead is in [Workflow](../Workflow/doc.md#add-a-dependency).

An unfilled `Deps` field is a nil func: it panics on first use, never silently. `verify` is
what stops that reaching a build.
