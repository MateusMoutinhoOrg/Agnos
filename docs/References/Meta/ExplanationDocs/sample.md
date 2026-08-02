# Dependency Mechanics

## Description
Explains how the library's dependency injection works: adapters build a `deps.Deps` struct of function fields, and the library calls those functions without knowing which adapter provided them.

---

## The Deps Contract

`deps.Deps` is a plain struct of function fields — not an interface. Every adapter fills in each field, and the library reaches dependencies only through it:

```go
type Deps struct {
    ExampleDepFunctionA func() int // each field is one injectable behavior
}
```

Because fields are plain functions, any single behavior can be replaced without writing a full adapter.

---

## Overwriting a Dependency

An adapter's output is just a struct value, so specific behaviors can be swapped before injection:

```go
myDeps := standard.New(3)

// Replace only the behavior you need to change
myDeps.ExampleDepFunctionA = func() int {
    return 404
}

l := lib.New(myDeps) // the library is unaware of the change
```
