# `sandbox/api/sandbox.go`

## `Sandbox`

Sandbox is the whole library: one field per contract declared in sandbox/api/, each built by the New<Contract> of its own package under sandbox/internal/. sandbox.New returns it, and nothing callable lives outside of it.

| Field | Type | Description |
| --- | --- | --- |
| `UserSandbox` | `UserSandbox` | UserSandbox is the part of the Sandbox the project declares itself, in sandbox/api/usersandbox.go: embedded, so each of its fields is read as sandbox.<Field> like any contract. |
| `Deps` | `*deps.Deps` | Deps is every capability the sandbox reaches the outside world through. It rides on the api so that a function handed the Sandbox holds the whole of what it needs, and can call another field of the api besides — which is what makes a field a caller replaced take effect everywhere. It is also the one field that does not cross into a consumer: an installed copy of this contract carries the api, never the wiring behind it. |
| `Actions` | `Actions` |  |
| `Cli` | `Cli` |  |
| `Config` | `Config` |  |

[every contract](doc.md)
