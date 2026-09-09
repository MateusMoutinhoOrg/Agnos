# DepList

What `{{.GeneratorName}} dep-install <dep>` renders into this project. A dep is named after the contract it
installs under `sandbox/deps/<dep>/`; the `Deps` field is the title-cased dir name, and the
adapter lib filling it lands under `adapters/libs/`. Signatures of the ones already installed
are in [PublicApi](../PublicApi/doc.md#dependency-contracts).

| Dep | `Deps` field | Adapter lib | Backed by | Provides |
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
| `sortdeps` | `Sortdeps` | `sortdeps` | `sort` | Sort a string slice, or any slice by a less function |
| `std` | `Std` | `std` | `time`, `fmt`, `runtime`, `os.Stdout/Stderr` | Clock, `Sprintf`, the host `Goos` and the three output channels. Installed by `cli-init` |
| `stringsdeps` | `Stringsdeps` | `stringsdeps` | `strings`, `strconv` | Text manipulation and string/number conversion |
| `templatedeps` | `Templatedeps` | `templatedeps` | `text/template` | Parse and execute one template over vars, with native funcs |

`{{.GeneratorName}} dep-list` prints the same names; `{{.GeneratorName}} dep-remove <dep>` takes one back out. Writing a
contract of your own instead is in
[Workflow](../Workflow/doc.md#add-a-dependency).

`adapters/availables/standard/new.go` is regenerated to bind every lib under `adapters/libs/`,
so an installed dep needs no wiring. An unfilled `Deps` field is a nil func: it panics on first
use, never silently.
