# DepList

What `agnos dep-install <dep>` renders into this project. A dep is named after the contract it
installs under `sandbox/deps/<dep>/`; the `Deps` field is the title-cased dir name, and the
adapter lib filling it lands under `adapters/libs/`. Signatures of the ones already installed
are in [PublicApi](../PublicApi/doc.md#dependency-contracts).

| Dep | `Deps` field | Adapter lib | Backed by | Provides |
|---|---|---|---|---|
| `argvdeps` | `Argvdeps` | `verb` | `github.com/MateusMoutinhoOrg/Verb` (pinned in `go.mod`) | Per-call argv parser. Installed by `cli-init` |
| `dbdeps` | `Dbdeps` | `keep` | `github.com/MateusMoutinhoOrg/Keep` (pinned) | Schema database |
| `embeddeps` | `Embeddeps` | `embeddeps` + `assets/asset.go` | `embed`, `text/template` | Read and render files compiled into the binary |
| `goimportsdeps` | `Goimportsdeps` | `goimportsdeps` | `go/parser` | Go source reader (package, imports, declarations) |
| `iodeps` | `Iodeps` | `iodeps` | `os`, `path/filepath` | Filesystem. `WriteFile` creates parents; `RemoveDir` removes files too |
| `requestdeps` | `Requestdeps` | `requestdeps` | `net/http` (30s timeout) | Per-call HTTP request |
| `rundeps` | `Rundeps` | `rundeps` | `os/exec` | Run a program to completion; stdout+stderr merged; non-zero exit is `Result.ExitCode`, not an error |
| `serializables` | `Serializables` | `serializables` | `gopkg.in/yaml.v3`, `encoding/json` | Generic JSON/YAML values |
| `std` | `Std` | `std` | `time`, `fmt`, `os.Stdout/Stderr` | Clock and the three output channels. Installed by `cli-init` |

`agnos dep-list` prints the same names; `agnos dep-remove <dep>` takes one back out. Writing a
contract of your own instead is in
[Workflow](../Workflow/doc.md#add-a-dependency).

`adapters/availables/standard/new.go` is regenerated to bind every lib under `adapters/libs/`,
so an installed dep needs no wiring. An unfilled `Deps` field is a nil func: it panics on first
use, never silently.
