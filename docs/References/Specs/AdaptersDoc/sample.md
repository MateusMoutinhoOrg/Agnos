# Adapters

## Description
Lists every adapter lib and assembly shipped under `adapters/`, and when to use each one. To build a new lib, follow [AddAdapterLib.md](/docs/Tutorials/AddAdapterLib.md).

---

## Available Adapters

| Adapter | Factory | Behavior | Use When |
|---------|---------|----------|----------|
| `availables/standard` | `standard.New` | Binds every lib under `adapters/libs/` | You want the default |
| `availables/frozen` | `frozen.New` | Binds the same libs, then replaces the clock with a fixed time | You need deterministic timestamps in tests |

This sample's `frozen` assembly is hypothetical, so its factory is unlinked: an entry links to its `docs/References/PublicApi/<pkg>.<Symbol>.md` page only when that page is real.

---

## Adapter Libs

| Lib | Fills | Backed by | Notes |
|-----|-------|-----------|-------|
| `libs/std` | `Deps.Std` | `time`, `fmt`, `os` | `Printf` to stdout; `Log` and `Error` to stderr |
| `libs/clockdeps` | `Deps.Clockdeps` | `time` | `Now` and `Sleep` |

---

## Embedded Libraries

`Deps` carries no whole third-party library in this sample; every field is a set of behaviors written in the lib itself.

---

## Standing Capabilities

`Clockdeps.Sleep` is declared and filled like any other dependency, but no action calls it yet. An assembly must fill it regardless: an unfilled field is a nil function the compiler does not catch, and it panics on first use.

| Field | Filled by | `standard` | `frozen` |
|-------|-----------|------------|----------|
| `Clockdeps` | `Bind` in `adapters/libs/clockdeps/Clock.go` | The real wall clock | A time fixed at construction |
