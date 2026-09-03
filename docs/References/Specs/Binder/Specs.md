# Binder Specification

## Description
Defines the required shape of a hand-written file of `sandbox/binds/` — a **binder**: the one function that assigns the sandbox's internal implementations onto the contract struct declared in the `sandbox/api/` file of the same name. `actions.go` is the one such file today; `cli.go` is generated. `verify` enforces the mirror and that the file declares only functions.

### Rules
- `package binds`. The file has the **same name** as the `api` file it fills, and there is exactly one binder per contract.
- One exported function `<Contract>Bind(deps *deps.Deps, sandbox *api.Sandbox)`, and nothing else at top level — no types, consts or vars (imports excepted).
- The function assigns **every** field of `sandbox.<Contract>`, in the order the contract declares them. Each assignment is a closure that forwards to the action package's public entry, passing `deps` first: `return buildAction.Build(deps, props)`.
- Action packages are imported under the alias `<name>Action` (`buildAction`, `depInstallAction`), so the closure body reads as the operation it forwards to.
- The binder contains no logic: no validation, no logging, no defaulting. Those belong in the action.
- `sandbox/new.go` — generated — calls every binder found in this directory; adding a file here is all the registration a new contract needs.

## Structure
1. **Package clause**: `package binds`.
2. **Imports**: `api`, `deps`, and one `<name>Action` alias per action forwarded to.
3. **The binder**: `func <Contract>Bind(deps *deps.Deps, sandbox *api.Sandbox)` with one assignment per contract field.

> **Note**: For a concrete example, refer to [sample.go](/docs/References/Specs/Binder/sample.go).
