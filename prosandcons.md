# Migração para interfaces: `sandbox/api` e `sandbox/deps`

## 1. Resumo

| | |
|---|---|
| Proposta | `sandbox/api/` e `sandbox/deps/` passam a conter **apenas interfaces Go**, no lugar de structs de campos-função |
| Ganho | completude checada pelo compilador, adapters com estado, implementação sem importar agnos, contrato portável |
| Custo | perde-se o override de um campo só, mocks parciais, e o shim remoto continua; ~1.500 call-sites + gêmeos em `assets/` |
| Inviável | "apenas interfaces" de forma literal: Props, Config, `CommandArg`… são **dados** e só virariam interface com getters |
| Recomendação | **híbrido**: interfaces em `sandbox/deps/<dep>/`, `sandbox/api/` continua em structs (seção 7) |

## 2. Estado atual

| Unidade | Forma hoje | Arquivo |
|---|---|---|
| `deps.Deps` | struct com 14 campos, um `<dep>.Sandbox` por dep | `sandbox/deps/deps.go` |
| `<dep>.Sandbox` | struct de campos `func` (`ReadFile func(path string) ([]byte, error)`) | `sandbox/deps/iodeps/iodeps.go` |
| Adapter | `Bind(deps *deps.Deps)` atribui closures | `adapters/libs/iodeps/iodeps.go` |
| Available | chama os 14 `Bind` | `adapters/availables/standard/new.go` |
| `api.Sandbox` | struct que embute `CliSandbox`, `UserSandbox` + `Deps`, `Config` | `sandbox/api/sandbox.go` |
| `api.Actions` | struct com ~81 campos `func(props XProps) error` | `sandbox/api/actions.go` |
| Props / Config / `Command*` | structs de dados puros | `sandbox/api/actions.go`, `command.go`, `config.go` |
| `CommandResponse` | struct de campos `func` (`SetStatus`, `Printf`…) | `sandbox/api/command.go` |
| Constructor | `sandbox.Actions = actions.NewActions(sandbox)`, "wrap, decore ou substitua" | `sandbox/constructors/actions/constructor.go` |
| Cobertura | `verify` exige que cada available preencha cada campo exatamente uma vez | `sandbox/internal/actions/verify/check_adapters.go` (`checkAvailableCoverage`) |
| Cópia remota | `sandbox/api` copiado para `sandbox/deps/<dep>/` do consumidor + shim gerado que converte struct a struct | `sandbox/internal/apishape/`, `check_api_shape.go`, `check_remote_deps.go` |

## 3. Modelo proposto

```go
// hoje                                        // proposto
type Sandbox struct {                          type Sandbox interface {
    ReadFile  func(path string) ([]byte, error)    ReadFile(path string) ([]byte, error)
    WriteFile func(path string, c []byte) error    WriteFile(path string, c []byte) error
}                                              }

// adapter hoje                                // adapter proposto
func Bind(deps *deps.Deps) {                   type Adapter struct{}
    deps.Iodeps = iodeps.Sandbox{              var _ iodeps.Sandbox = (*Adapter)(nil)
        ReadFile: func(p string) ... {...},    func (a *Adapter) ReadFile(p string) ... {...}
    }                                          func Bind(deps *deps.Deps) { deps.Iodeps = &Adapter{} }
}
```

O que **não** vira interface sem getters: `BuildProps`, `StartProps`, todos os `*Props`, `Config`,
`CommandArg`, `CommandFlag`, `CommandFailure`, `ServerProps`, constantes (`ExitOk`, `RuntimeGo`,
`ArgType`, `FlagType`). Um pacote "só de interfaces" ainda precisa de um pacote irmão de dados.

## 4. Prós

| # | Pró | Impacto | Onde |
|---|---|---|---|
| P1 | Completude checada pelo compilador: um método faltando não compila | some o panic de campo nil em runtime; `checkAvailableCoverage` pode ser simplificado | `check_adapters.go` |
| P2 | Satisfação implícita: qualquer tipo com os métodos implementa o contrato, sem importar agnos | lib de terceiros ou outro repo vira adapter direto | `adapters/libs/` |
| P3 | Adapter com estado (conexão, cache, pool) como campos de struct, não capturados por closure | código de adapter mais legível e testável | `adapters/libs/*` |
| P4 | Idiomático em Go; gopls "go to implementation", linters e docs padrão funcionam | menos convenção própria a aprender | todo o repo |
| P5 | Um único nil-check por dep (a interface inteira), não por função | erro de wiring aparece cedo | `cmd/main/` |
| P6 | Contrato separado da implementação de forma explícita: portar para outro runtime (WASM, mock, remoto) é implementar a interface | portabilidade | `sandbox/deps/` |
| P7 | Interface cujos métodos usam só builtins é idêntica entre cópias: atravessa sem shim | menos código gerado em `apishape` para deps simples | `sandbox/internal/apishape/` |
| P8 | Valor zero explícito: `Deps{}` com interfaces nil é obviamente vazio; struct de funções parece preenchida | diagnóstico mais claro | `sandbox/deps/deps.go` |

## 5. Contras

| # | Contra | Impacto | Onde |
|---|---|---|---|
| C1 | **Perde o override de um campo só.** Hoje `sandbox.Actions.Build = wrap(sandbox.Actions.Build)` decora uma action; com interface é preciso um tipo wrapper com todos os métodos (~81 em `Actions`), ou embutir a interface e sobrescrever um | quebra a promessa do constructor ("wrap, decore ou substitua") e o comentário de `api.Sandbox` ("um campo substituído vale em todo lugar") | `sandbox/constructors/*`, `sandbox/api/sandbox.go` |
| C2 | Mocks parciais viram fakes completos; embutir a interface para sobrescrever um método dá panic em nil no resto | testes de exemplo e adapters de teste mais longos | `examples/`, `adapters/libs/*` |
| C3 | "Apenas interfaces" é inalcançável: dados exigiriam getters/setters | explosão de código e tokens, contra "docs/código densos" do `CLAUDE.md` | `sandbox/api/actions.go`, `command.go`, `config.go` |
| C4 | Shim remoto continua: um método que recebe/retorna um tipo do próprio pacote (`BuildProps`, `ServerProps`) não é idêntico entre a cópia e o original | `apishape` precisa gerar wrappers de método em vez de closures, sem remover a complexidade | `sandbox/internal/apishape/converters.go`, `convert_expression.go` |
| C5 | Embedding de partes (`CliSandbox`, `UserSandbox`, `<x>config.go`) com interfaces: nomes de método colidem, e embutir interface em struct só repassa, não compõe | o modelo "cada mecânica adiciona uma parte" fica mais frágil | `sandbox/api/*sandbox.go`, `apishape` |
| C6 | Volume da migração: ~1.460 usos de `sandbox.Deps.` e ~80 de `sandbox.Actions.` fora de `assets/`, cada um com gêmeo em `assets/` | migração grande, só viável por template + bootstrap | `sandbox/`, `assets/` |
| C7 | `assets/deplist/<dep>/**` deve renderizar byte a byte igual à cópia usada aqui | todo dep do catálogo e todo adapter reescritos em par | `assets/deplist/`, `assets/adapterlist/` |
| C8 | `verify` reescrito: `check_contracts`, `check_adapters`, `check_api_shape`, `check_remote_deps` leem forma de struct | risco de regressão nas checagens | `sandbox/internal/actions/verify/` |
| C9 | **Breaking change** para todo projeto gerado e para toda dep remota já copiada (shape antigo) | exige versão major e caminho de migração | projetos externos |
| C10 | Docs: `docs/PublicApi/` é gerado dos comentários de campo; comentário passa a ser de método. `Contributing`, `Rules`, `Structure`, `LibUsage` descrevem o padrão atual | reescrita dos templates em `assets/doc/` | `assets/doc/docs/` |
| C11 | Todos os goldens de `examples/` que contêm código/docs gerados mudam | `exec-test --update` amplo, difícil de revisar | `examples/` |
| C12 | Chamada via interface não é inlinada (custo desprezível na prática) | irrelevante para cli; marginal em server | — |
| C13 | Campos de função são valores: dá para trocar em runtime, compor por campo, montar um available misturando funções de adapters diferentes. Interface é tudo-ou-nada por dep | um available não pode mais pegar `ReadFile` de um adapter e `WriteFile` de outro sem wrapper | `adapters/availables/` |

## 6. Matriz por unidade

| Unidade | Recomendação | Motivo |
|---|---|---|
| `<dep>.Sandbox` (`sandbox/deps/<dep>/`) | **migrar** | maior ganho (P1, P2, P3, P6); override por campo raramente usado em deps |
| `deps.Deps` | manter struct, campos passam a ser interfaces | é agregador, não contrato |
| `api.Actions` | **manter struct de funções** | é onde C1 dói: decorar uma action é o caso de uso do constructor |
| `api.Cli`, `api.Server`, partes `<x>sandbox.go` | manter struct | dependem de embedding (C5) |
| Props, Config, `Command*`, `CommandFailure` | manter dados | C3 |
| `CommandResponse`, `serverdeps.Response` | avaliar caso a caso | interface ajudaria fakes completos, mas hoje o dispatch monta closures por chamada |

## 7. Alternativa híbrida recomendada

1. `sandbox/deps/<dep>/` declara uma **interface** `Sandbox`; adapters viram tipos com métodos e
   `var _ <dep>.Sandbox = (*Adapter)(nil)`.
2. Para manter o override pontual, cada dep pode ganhar um gerado `Funcs` (struct de campos-função
   que implementa a interface), servindo de adaptador nos dois sentidos — decorar continua sendo
   trocar um campo.
3. `sandbox/api/` continua em structs; o shim remoto e o modelo de partes não mudam.
4. Opcional: interfaces de leitura menores (`Reader`, `Writer`) para quem precisa só de parte da dep
   — o ganho real de portabilidade em Go vem de interfaces pequenas.

## 8. Impacto e ordem da migração (híbrida)

| Passo | O que muda | Onde |
|---|---|---|
| 1 | contrato de cada dep vira interface | `assets/deplist/<dep>/sandbox/deps/<dep>/`, depois `sandbox/deps/<dep>/` |
| 2 | adapters viram tipos com métodos | `assets/adapterlist/<adapter>/`, `adapters/libs/*` |
| 3 | template de nova dep/adapter (`add-dep`, `add-adapter`) | `sandbox/internal/actions/add_dep/`, templates em `assets/` |
| 4 | `verify`: cobertura vira checagem de compilação; `check_contracts` lê métodos | `sandbox/internal/actions/verify/` |
| 5 | geração de `docs/PublicApi/` a partir de comentários de método | collector do PublicApi, `assets/doc/` |
| 6 | docs de regra e contribuição | `assets/doc/docs/Rules/doc.md`, `docs/Contributing/doc.md`, `CLAUDE.md` |
| 7 | bootstrap + idempotência + goldens | `go build -o release/bootstrap.bin ./cmd/main && ./release/bootstrap.bin build`, `exec-test --update` |
| 8 | versão major + nota de migração para projetos existentes | `AgnosConfig/project.yaml` |

Call-sites `sandbox.Deps.<Dep>.<Fn>(...)` **não mudam** de sintaxe: chamar método e chamar campo
função se escrevem igual. O custo está em declarações, adapters, verify, apishape e docs — não nos
~1.460 usos.

## 9. Riscos e perguntas abertas

- `apishape`: vale manter cópia de `sandbox/api` + shim, ou interfaces em `deps` permitem que um
  repo remoto exponha diretamente uma interface só de builtins?
- Algum projeto (ou o próprio agnos) troca um campo de dep em runtime? Se sim, C13 vira bloqueio
  sem o `Funcs` da seção 7.
- `interviewer` e `serverdeps` retornam valores com estado (`Server`); confirmar que o shape
  interface → interface atravessa o shim.
- Migrar de uma vez ou por dep? Por dep exige que `verify` aceite os dois shapes durante a
  transição.
