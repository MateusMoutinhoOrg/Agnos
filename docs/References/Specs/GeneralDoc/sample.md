# Configuration Loading

## Description
How a generated project reads its configuration at build time and where the defaults come from.

---

## Files

### `/AgnosConfig/`

| File | Description |
|------|-------------|
| `project.yaml` | `name`, `version`, `description` of the project |
| `ignore.yaml` | Paths hidden from every listing |

---

## Loading Order

Configuration is resolved from the lowest to the highest precedence:

1. The defaults declared in `projectconf.NewEmpty`.
2. The values found in `project.yaml`.
3. The flags passed on the command line.

---

## Usage

```go
// loadProjectConf reads project.yaml through the transaction and fails fast
// when it is missing: `agnos start` is a prerequisite for `agnos build`.
conf, err := loadProjectConf(deps, io, path)
if err != nil {
    return err
}
```

Actions receive the resolved values as plain arguments — see [AddAction.md](/docs/Tutorials/AddAction.md).
