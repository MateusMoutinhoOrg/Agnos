# Configuration Loading

## Description
How the project reads its configuration at startup and where the defaults come from.

---

## Files

### `/config/`

| File | Description |
|------|-------------|
| `config.go` | Loads and validates environment configuration |
| `defaults.go` | Defines default values for all settings |

---

## Loading Order

Configuration is resolved from the lowest to the highest precedence:

1. The defaults declared in `defaults.go`.
2. The values found in the environment.
3. The flags passed on the command line.

---

## Usage

```go
// Load resolves the configuration and fails fast on an invalid value.
cfg, err := config.Load()
if err != nil {
    return err
}
```

Adapters receive the resolved values as plain arguments — see [DepsMechanic.md](/docs/Explanations/DepsMechanic.md).
