# Runtimes

## Description
Every runtime `agnos build` and `agnos verify` can hand a project to, with the steps each one runs.

---

## Runtimes

| Runtime | Steps | Use When |
|---------|-------|----------|
| `go` | `go mod tidy`, then `go build` over the schema dirs that exist | The default: success means the toolchain accepted the tree |
| `none` | Nothing | No toolchain on the machine, or hand-written code known not to compile yet |

---

## Constants

### [api.RuntimeGo / api.RuntimeNone](/docs/PublicApi/api.BuildProps/doc.md)
The two accepted values of `api.BuildProps.Runtime`; any other string is a usage error.
