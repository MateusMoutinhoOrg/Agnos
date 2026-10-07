# Plano de ação: deps como interfaces (alternativa híbrida)

Base: [interfacechange.md](interfacechange.md), seção 7.

## 1. Escopo

| Entra | Fica de fora |
|---|---|
| O `Sandbox` de cada dep do **catálogo** (`assets/deplist/<dep>/`) vira `interface` | `sandbox/api/*` (Actions, Cli, Props, Config, Command*) |
| Um `Funcs` gerado por dep, para continuar trocando uma função só | deps **remotas** (cópia de outro `sandbox/api`, que continua struct) |
| Adapters do catálogo viram tipos com métodos | tipos aninhados de função: `serverdeps.Server/Request/Response`, parser do `argvdeps`, `requestdeps` (fase futura, seção 6) |
| `verify`, docs, release | `apishape` e o shim remoto |

## 2. Desenho final

```go
// sandbox/deps/iodeps/iodeps.go — escrito à mão, vem de assets/deplist/iodeps/
// Sandbox is the filesystem library injected whole as the Deps.Iodeps field.
type Sandbox interface {
	// ReadFile returns the whole content of the file at path.
	ReadFile(path string) ([]byte, error)
	// WriteFile writes content to path, creating any missing parent directory.
	WriteFile(path string, content []byte) error
	// ...
}
```

```go
// sandbox/deps/iodeps/funcs.go — GERADO por build a partir da interface
// Funcs is Sandbox spelled as function fields, so one of them can be replaced.
type Funcs struct {
	ReadFile  func(path string) ([]byte, error)
	WriteFile func(path string, content []byte) error
}

// Wrap returns a Sandbox whose every method calls the matching field of funcs.
func Wrap(funcs Funcs) Sandbox { return wrapped{funcs} }

// Unwrap returns the methods of sandbox as function fields.
func Unwrap(sandbox Sandbox) Funcs {
	return Funcs{ReadFile: sandbox.ReadFile, WriteFile: sandbox.WriteFile}
}

type wrapped struct{ funcs Funcs }

func (w wrapped) ReadFile(path string) ([]byte, error) { return w.funcs.ReadFile(path) }
func (w wrapped) WriteFile(path string, content []byte) error { return w.funcs.WriteFile(path, content) }
```

Go não deixa um tipo ter campo e método com o mesmo nome, por isso são dois tipos (`Funcs` e `wrapped`).

```go
// adapters/libs/iodeps/iodeps.go — vem de assets/adapterlist/iodeps/
type Adapter struct{}

var _ iodeps.Sandbox = (*Adapter)(nil) // completude checada pelo compilador

// Bind fills deps.Deps.Iodeps. Assinatura inalterada.
func Bind(deps *deps.Deps) { deps.Iodeps = &Adapter{} }

func (adapter *Adapter) ReadFile(path string) ([]byte, error) { return os.ReadFile(path) }
```

```go
// troca pontual, ex. --quiet em sandbox/internal/commands/project/InternalPureHandler.go
funcs := std.Unwrap(sandbox.Deps.Std)
funcs.Log = func(format string, a ...any) (int, error) { return 0, nil }
sandbox.Deps.Std = std.Wrap(funcs)
```

Chamar continua igual: `sandbox.Deps.Iodeps.ReadFile(path)`.

## 3. Fases

Cada fase termina verde: `go build -o release/bootstrap.bin ./cmd/main && ./release/bootstrap.bin build`,
`./release/bootstrap.bin build -q && git diff --quiet`, `exec-test --only <exemplos tocados>`.

| Fase | O que fazer | Arquivos |
|---|---|---|
| **F0** Preparação | branch `deps-interfaces`; escrever a regra nova ("o `Sandbox` de uma dep do catálogo é interface; `funcs.go` é gerado") | `assets/doc/docs/Rules/doc.md` |
| **F1** Gerador de `Funcs` | collector que lê cada `sandbox/deps/<dep>/` com `Deps.Goimportsdeps.Parse` (já devolve `Kind == "interface"` e `Methods`) e, quando `Sandbox` é interface, renderiza o template novo → `sandbox/deps/<dep>/funcs.go`. Sem efeito enquanto nenhuma dep for interface | novo `sandbox/internal/actions/build/collect_dep_funcs.go` (molde: `collect_deps_api.go`), novo `assets/templates/dep_funcs.go`, ligação em `build_internal.go`, linha em `assets/doc/docs/GeneratedFiles/doc.md` |
| **F2** Piloto `std` | contrato → interface; adapter → `type Adapter` + `var _`; `--quiet` passa a usar `Unwrap/Wrap` | `assets/deplist/std/sandbox/deps/std/std.go` + `sandbox/deps/std/std.go`; `assets/adapterlist/std/adapters/libs/std/std.go` + `adapters/libs/std/std.go`; `sandbox/internal/commands/project/InternalPureHandler.go` |
| **F3** Demais deps do catálogo | mesmo molde da F2, **uma dep por commit**, das menores para as maiores: `rundeps`, `templatedeps`, `sortdeps`, `embeddeps`, `goimportsdeps`, `interviewer`, `reflectdeps`, `hashdeps`, `iodeps`, `argvdeps`*, `serverdeps`*, `stringsdeps`, `serializables`, depois as não instaladas aqui (`envdeps`, `jwtdeps`, `passworddeps`, `randdeps`, `ratelimitdeps`, `requestdeps`*, `signaldeps`, `timedeps`) | `assets/deplist/<dep>/`, `assets/adapterlist/<adapter>/` (incl. `reflectsort`), cópias instaladas em `sandbox/deps/` e `adapters/libs/` |
| **F4** verify | mensagens "nil func" → "nil interface" (a cobertura por available continua: dois `Bind` no mesmo campo ainda se sobrescrevem em silêncio); checagem nova: dep de origem catálogo declara `Sandbox` interface e tem `funcs.go`; dep remota pode ser struct | `sandbox/internal/actions/verify/check_adapters.go`, `check_contracts.go` ou `check_deplist.go` |
| **F5** Docs | "Custom deps" passa a mostrar `Unwrap/Wrap`; receita de dep/adapter; descrição da camada deps | `assets/doc/docs/LibUsage/doc.md`, `assets/doc/docs/Structure/doc.md`, `docs/Contributing/doc.md`, `CLAUDE.md` (seções "deps layer" e "Naming is load-bearing") |
| **F6** Release | versão major; regenerar goldens; nota de migração para projetos gerados | `AgnosConfig/project.yaml`, `exec-test --update`, `examples/**/result.yaml` |

\* Na F3, estas mudam só o `Sandbox` de topo; `NewServer`/`New`/`NewRequest` continuam devolvendo os
structs aninhados atuais.

## 4. O que muda

| Área | Mudança | Fase |
|---|---|---|
| `assets/deplist/<dep>/sandbox/deps/<dep>/*.go` + cópia instalada | `type Sandbox struct{ F func(...) }` → `type Sandbox interface{ F(...) }`, comentários de campo → de método | F2–F3 |
| `sandbox/deps/<dep>/funcs.go` | **novo, gerado** a cada build | F1 |
| `assets/templates/dep_funcs.go` | **novo** template | F1 |
| `sandbox/internal/actions/build/` | collector + render do `funcs.go` | F1 |
| `assets/adapterlist/<adapter>/adapters/libs/<adapter>/*.go` + cópia instalada | closures → métodos de `Adapter`; `var _ <dep>.Sandbox = (*Adapter)(nil)` | F2–F3 |
| `sandbox/internal/commands/project/InternalPureHandler.go` | `Std.Log = ...` → `Unwrap/Wrap` | F2 |
| `sandbox/internal/actions/verify/` | mensagens + checagem de forma da dep | F4 |
| `assets/doc/docs/{Rules,LibUsage,Structure,GeneratedFiles}/doc.md`, `docs/Contributing/doc.md`, `CLAUDE.md` | descrevem o novo padrão | F0, F1, F5 |
| `docs/PublicApi/` | renderiza interfaces (o `collect_public_api.go` já tem `IsInterface`/`Methods`; conferir a página) | F2 |
| `examples/**` | goldens com `docs/` e adapters gerados | F6 |

## 5. O que não muda

| Item | Por quê |
|---|---|
| `assets/sandbox-deps/sandbox/deps/deps.go` | o campo continua `{{.Title}} {{.Name}}.Sandbox`, seja interface ou struct |
| ~1.460 chamadas `sandbox.Deps.<Dep>.<Fn>(...)` | chamar método e chamar campo-função se escrevem igual |
| `Bind(deps *deps.Deps)` e `assets/templates/available_new.go` | o binder só passa a atribuir a dep inteira |
| `sandbox/api/*`, `sandbox/internal/apishape/`, `assets/templates/remote_shim.go`, `add_dep/remote_dep.go` | `sandbox/api` continua struct, então a dep remota copiada também |
| `checkAvailableCoverage` (lógica) | continua necessária contra dois adapters no mesmo campo |

## 6. Fase futura (fora deste plano)

`serverdeps.Server/Request/Response`, o parser do `argvdeps` e o `requestdeps` como interfaces. Eles
aparecem em `sandbox/api` (`api.Server`, `response *serverdeps.Response` dos handlers) e nos
`handle_*.go` que já são do projeto, então mudá-los quebra código do usuário. Decidir só depois da F6.

## 7. Riscos

| Risco | Mitigação |
|---|---|
| Regra "`assets/deplist/<dep>/**` renderiza byte a byte" vs `funcs.go` gerado | `funcs.go` fica fora do deplist e da comparação de `check_deplist.go`; conferir na F1 |
| Projeto gerado com adapter escrito à mão quebra ao atualizar a dep | nota de migração: `deps.X = x.Wrap(x.Funcs{...})` aceita as closures antigas sem reescrever |
| `Wrap(Unwrap(a))` perde o tipo concreto (`a.(*Adapter)` passa a falhar) | documentar em LibUsage; nenhum uso no repo hoje |
| Goldens mudam em massa e a revisão fica difícil | `exec-test --update` só na F6; nas fases anteriores, `--only` |
| Dep com interface e dep remota com struct convivem em `Deps` | é o comportamento desejado; a checagem da F4 aceita os dois por origem |

## 8. Critérios de pronto

- `verify` limpo; `build` compila e é idempotente; `exec-test` verde.
- Toda dep do catálogo declara `Sandbox` interface e tem `funcs.go` gerado.
- Todo adapter do catálogo tem `var _ <dep>.Sandbox = (*Adapter)(nil)`.
- `agnos build -q` continua silenciando o `Log`.
- `docs/PublicApi/` lista os métodos de cada dep com descrição.
