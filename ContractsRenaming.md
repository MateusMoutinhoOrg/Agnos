### Objective
Rename the contracts dir and objects with names more easly to understand for the user

## Dirs:
| Dir  | New Dir | Comment |
| ---- | -------- | ------- |
| `sandbox/contracts/api` | `sandbox/contracts/lib` | O pacote descreve a lib construida por `sandbox.New`, nao uma API externa. O codigo ja usava esse vocabulario nos comentarios |

## Files:
| File | New File | Comment |
| ---- | -------- | ------- |
| `sandbox/contracts/api/api.go` | `sandbox/contracts/lib/sandbox.go` | O arquivo passa a ter o nome da struct que ele declara |

### Classes:

| Class | Class Location |  New Class | New Location | Comment |
| ----- | ---------      |  --------  | ----------  | ------- |
| api.SandBox | contracts/api/api.go | lib.SandBox  | contracts/lib/sandbox.go | So o pacote muda. `lib.Lib` seria stutter e `SandBox` descreve melhor o ambiente isolado |
| api.CliApi | contracts/api/cli.go | lib.CliApi | contracts/lib/cli.go | Idem |
| api.CoreApi | contracts/api/core.go | lib.CoreApi | contracts/lib/core.go | Idem |
| api.Config | contracts/api/config.go | lib.Config | contracts/lib/config.go | Idem |

### Fora de escopo

- `docs/References/Specs/` — desatualizado, nao acompanha o rename.
- `adapters/standard/keep.go` — o import `keepapi` aponta para o pacote `contracts/api` do repositorio Keep, que nao muda.
- `sandbox/contracts/deps/{keepdeps,verbdeps,iodeps,requestdeps,embeddeps}` — os comentarios "mirrors the embedded library's api.X" citam as libs externas.
