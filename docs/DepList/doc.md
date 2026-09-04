# DepList

What `agnos dep-install <dep>` renders. A dep is named after the contract it installs under `sandbox/deps/<dep>/`; the `Deps` field is the title-cased dir name. Contract signatures are in [PublicApi](/docs/PublicApi/doc.md#dependency-contracts).

| Dep | `Deps` field | Adapter lib | Backed by | Provides |
|---|---|---|---|---|
| `argvdeps` | `Argvdeps` | `adapters/libs/verb` | `github.com/MateusMoutinhoOrg/Verb` (pinned in `go.mod`) | Per-call argv parser. Installed by `cli-init` |
| `dbdeps` | `Dbdeps` | `adapters/libs/keep` | `github.com/MateusMoutinhoOrg/Keep` (pinned) | Schema database. Not installed in Agnos itself |
| `embeddeps` | `Embeddeps` | `adapters/libs/embeddeps` + `assets/asset.go` | `embed`, `text/template` | Read and render files compiled into the binary |
| `goimportsdeps` | `Goimportsdeps` | `adapters/libs/goimportsdeps` | `go/parser` | Go source reader (package, imports, declarations) |
| `iodeps` | `Iodeps` | `adapters/libs/iodeps` | `os`, `path/filepath` | Filesystem. `WriteFile` creates parents; `RemoveDir` removes files too |
| `requestdeps` | `Requestdeps` | `adapters/libs/requestdeps` | `net/http` (30s timeout) | Per-call HTTP request |
| `rundeps` | `Rundeps` | `adapters/libs/rundeps` | `os/exec` | Run a program to completion; stdout+stderr merged; non-zero exit is `Result.ExitCode`, not an error |
| `serializables` | `Serializables` | `adapters/libs/serializables` | `gopkg.in/yaml.v3`, `encoding/json` | Generic JSON/YAML values |
| `std` | `Std` | `adapters/libs/std` | `time`, `fmt`, `os.Stdout/Stderr` | Clock and the three output channels. Installed by `cli-init` |

`adapters/availables/standard/new.go` is generated to bind every lib under `adapters/libs/`. An unfilled `Deps` field panics on first use.
