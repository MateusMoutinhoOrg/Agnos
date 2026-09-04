# The Quiet Flag

## Description
Explains how `--quiet` works everywhere without a handler doing anything: the generated dispatch swaps one function on the injected `Deps`.

---

## Three Channels

`deps.Std` carries three writers with one meaning each — `Printf` for the result, `Log` for progress, `Error` for failures. Only progress is noise a script wants gone:

```go
deps.Std.Log("build started with path %s \n", path) // progress: silenced by --quiet
deps.Std.Printf("%s\n", dep)                        // result: never silenced
```

---

## One Reassignment

When a command declares a boolean flag named `quiet`, its dispatch reads it first and, if present, replaces the progress writer before the handler runs:

```go
// generated into climain.go
if verb.IsPresent([]string{"--quiet", "-q"}) {
	silenceLogs(deps) // deps.Std.Log = a no-op
}
```

Because `Deps` is a struct of function fields, the swap is an assignment, and every action downstream is silenced with it — see [StructContracts](/docs/StructContracts/doc.md).
