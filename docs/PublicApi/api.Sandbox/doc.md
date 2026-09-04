# `api.Sandbox`

**Type:** Struct

## Definition

```go
type Sandbox struct {
	Actions Actions
	Cli     Cli
}
```

## Description

The library entry point, returned by [`sandbox.New`](/docs/PublicApi/sandbox.New/doc.md). It is **generated** by `agnos build` with one field per file of `sandbox/api/` other than itself — the collector `CollectConstructors` lists the directory and title-cases the names — so a project that adds an `api/<thing>.go` contract and a `binds/<thing>.go` binder gets a `Thing` field on the next build. Each field is a struct of function fields filled by the matching binder; see [StructContracts](/docs/StructContracts/doc.md).

## Fields

| Field | Description |
| :--- | :--- |
| `Actions` [`api.Actions`](/docs/PublicApi/api.Actions/doc.md) | Every operation `agnos` performs, callable from Go with the same semantics as the commands. |
| `Cli` [`api.Cli`](/docs/PublicApi/api.Cli/doc.md) | The generated command-line interface: parsing, dispatch, usage errors, exit code. |

## Examples

```go
lib := agnoslib.New(&deps)

// The library way: a typed call, an error back.
if err := lib.Actions.Verify("./my-tool"); err != nil {
	log.Fatal(err)
}

// The CLI way: an argument vector in, an exit code out.
code := lib.Cli.CliMain([]string{"verify", "--path", "./my-tool"})
```
