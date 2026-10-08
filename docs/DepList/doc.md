# DepList

The catalogue `agnos add-dep <dep>` installs from. A **dep** is the contract, one
directory under `sandbox/deps/<dep>/`, filling the `Deps` field of the same name title-cased,
a trailing `deps` spelled as its own word (`argvdeps` fills `ArgvDeps`).
An **adapter** is one implementation of it, one directory under `adapters/impls/<adapter>/`.
The two are separate catalogues because one contract may have several implementations — how
that works is [Adapters](../Adapters/doc.md). Signatures of the contracts already installed are
in [PublicApi](../PublicApi/doc.md#dependency-contracts).

| Dep | `Deps` field | Adapters | Backed by | Provides |
|---|---|---|---|---|
| `OpinionatedAgnosCli` | `OpinionatedAgnosCli` | `OpinionatedAgnosCli` | `strings`, `strconv`, `regexp`, `reflect` | The cli mechanic: the command types `api/{cli,command,trigger}.go` alias, the dispatch chain, binding, failures, triggers. Installed by `cli-init` |
| `OpinionatedAgnosDatabase` | `OpinionatedAgnosDatabase` | `OpinionatedAgnosDatabase` | the `databasedeps` dep | The readers and filters every generated `methods.go` shares. Installed by `database-init` |
| `OpinionatedAgnosFront` | `OpinionatedAgnosFront` | `OpinionatedAgnosFront` | `strings`, the `embeddeps` dep | The file layer of `assets/front/`: which file a path names, kept inside the tree, and its media type. Installed by `front-init` |
| `OpinionatedAgnosServer` | `OpinionatedAgnosServer` | `OpinionatedAgnosServer` | `strings`, `strconv`, `regexp`, `sort`, `reflect` | The server mechanic: the route types `api/{server,route}.go` alias, the request chain, binding, json-schema, the writers. Installed by `server-init` |
| `argvdeps` | `ArgvDeps` | `stdargv` | `strings`, `strconv`, `time` | Per-call argv parser. Installed by `cli-init` |
| `embeddeps` | `EmbedDeps` | `goembed` + `assets/asset.go` | `embed`, `text/template` | Read and render files compiled into the binary |
| `envdeps` | `EnvDeps` | `osenv` | `os` | Read an environment variable — how a secret reaches the process without travelling in argv or a file. Installed by `backoffice-init` |
| `goimportsdeps` | `GoimportsDeps` | `stdgoimports` | `go/parser` | Go source reader (package, imports, declarations) |
| `hashdeps` | `HashDeps` | `sha256hash` | `crypto/sha256`, `encoding/hex` | SHA-256 of a byte slice, lower-case hex |
| `interviewdeps` | `InterviewDeps` | `ttyinterview` | `os`, `os/exec`, `bufio` | Ask a person a question and get the answer back typed. Arrow-key menus over stdin in raw mode, numbered prompts where there is no terminal |
| `iodeps` | `IoDeps` | `osio` | `os`, `path/filepath` | Filesystem. `WriteFile` creates parents; `RemoveDir` removes files too; `Join`/`Dir` build host paths |
| `jwtdeps` | `JwtDeps` | `golangjwt` | `github.com/golang-jwt/jwt/v5` | Sign and check HS256 JSON Web Tokens; claims are flat builtins (`Id`, `Subject`, `IssuedAt`, `ExpiresAt`, `Ip`). Installed by `backoffice-init` |
| `passworddeps` | `PasswordDeps` | `pbkdf2password` | `crypto/pbkdf2` | Salted PBKDF2-HMAC-SHA256 password hashes (600k iterations) and their constant-time check. Installed by `backoffice-init` |
| `randdeps` | `RandDeps` | `cryptorand` | `crypto/rand`, `encoding/hex` | Cryptographically secure random bytes, lower-case hex. Installed by `backoffice-init` |
| `ratelimitdeps` | `RatelimitDeps` | `memoryratelimit` | `sync`, `time` | In-memory fixed-window counters by key, safe across concurrent requests. Installed by `backoffice-init` |
| `reflectdeps` | `ReflectDeps` | `stdreflect` | `reflect` | Build, fill and call a value whose type is known only at run time |
| `requestdeps` | `RequestDeps` | `nethttprequest` | `net/http` (30s timeout) | Per-call HTTP request |
| `rundeps` | `RunDeps` | `osexecrun` | `os/exec` | Run a program to completion; stdout+stderr merged; non-zero exit is `Result.ExitCode`, not an error |
| `serializabledeps` | `SerializableDeps` | `stdserializable` | `encoding/json` + a bundled YAML codec | Generic JSON/YAML values. The YAML side reads the block subset (no anchors, aliases or explicit tags). Installed by `server-init` |
| `serverdeps` | `ServerDeps` | `nethttpserver` | `net/http`, `net/url`, `net` | Http server: opens the port, applies timeouts, hands every request to one handler; headers, cookies, a form body, the connection's client ip (`GetClientIp`, also answered as `X-Client-Ip`). Installed by `server-init` |
| `signaldeps` | `SignalDeps` | `ossignal` | `os/signal` | Run a function once when the process is asked to stop. Installed by `server-init`, for a graceful shutdown |
| `sortdeps` | `SortDeps` | `stdsort`, `reflectsort` | `sort` / `reflect` | Sort a string slice, or any slice by a less function |
| `stddeps` | `StdDeps` | `osstd` | `time`, `fmt`, `runtime`, `os.Stdout/Stderr` | Clock, `Sprintf`, the host `GOOS` and the three output channels. Installed by `cli-init` |
| `stringsdeps` | `StringsDeps` | `stdstrings` | `strings`, `strconv` | Text manipulation and string/number conversion. Installed by `cli-init` |
| `templatedeps` | `TemplateDeps` | `texttemplate` | `text/template` | Parse and execute one template over vars, with native funcs |
| `timedeps` | `TimeDeps` | `stdtime` | `time` | A Unix instant to and from a UTC date spelled by a Go layout. Installed by `backoffice-init` |

The first adapter of each row is the dep's `default-adapter`, the one `add-dep` installs when
`--adapter` names no other. The four `OpinionatedAgnos<X>` rows are the opinionated libs: each carries
a mechanic rather than a library's raw capability ([Rules](../Rules/doc.md#layers)).

```bash
agnos list-deps                       # this table, plus what is installed here
agnos list-adapters                   # one row per adapter, and which binding binds it
agnos add-dep sortdeps                # contract + default-adapter
agnos add-dep sortdeps --adapter reflectsort
agnos remove-dep sortdeps             # refused while an adapter fills it
agnos remove-dep sortdeps --with-adapters
```

Writing a contract of your own instead is in [Workflow](../Workflow/doc.md#add-a-dependency).

An unfilled `Deps` field is a nil func: it panics on first use, never silently. `verify` is
what stops that reaching a build.
