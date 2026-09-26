# Server — proposta de mudanças

**Status (2026-09-26): implementado**, exceto os itens marcados como pendentes abaixo. Onde a
implementação divergiu do plano:

| Item | Plano | Como ficou | Por quê |
|---|---|---|---|
| §1.1 `prefix` | por segmento sempre | por segmento **em path**; em valor de parâmetro é prefixo de texto | `Bearer abc` não tem segmentos |
| §5.2 binding | só tipados sem barra | todo path de **um segmento** (`start == end`) entrega o segmento puro (`"hello"`); range mantém `"/a/b"` | `{slug}` chegava como `"/hello"` |
| §3.5 405 | "nenhuma rota rodou" | rotas `ANY` não contam | um guard `ANY` escondia todo 405 |
| §3.4 middleware | trigger padrão `/<name>` | trigger padrão `/`, categoria `Middleware` | middleware sem trigger guarda tudo |
| §7.2 `add-route` | `--help`/`--category` opcionais | opcionais; categoria padrão `Routes` | — |
| §6.1 form | `Request.ReadForm` | idem, reusando o body em cache | `net/url` fica no adapter |
| §9 aviso "guard só GET" | aviso do verify | **pendente** — o verify só tem erros, não avisos | — |


Escopo: `sandbox-server` (templates em `assets/sandbox-server/`, `assets/doc-server/`), os
comandos `*-route`, `*-path`, `*-parameter`, `*-body*`, a parsable `routeconf`, o collector
`collect_routes.go` e o `check_routes.go`. Tudo abaixo é implementado no template e rodado via
bootstrap; nenhum arquivo gerado é editado à mão.

## 0. Diagnóstico (medido num projeto novo, `start` + `server-init`)

| Caso | Hoje | Esperado | Causa |
|---|---|---|---|
| `add-route admin --trigger /admin --trigger-type starts-with` | `unknown trigger type` | aceitar | só existe o nome `prefix` |
| `prefix /admin` × `GET /administrator`, `/admin-x` | **200** | 404 | `MatchTrigger` usa `HasPrefix` no texto, ignora segmentos |
| `GET /ADMIN` | 404 | configurável | comparação sempre case-sensitive |
| guard `/admin` criado depois de `admin` | roda **depois** | antes | toda rota nasce `priority: 0`, empate por nome, negativo proibido |
| guard para todos os métodos | 7× `-m` | `-m ANY` | `ANY`/`*` recusados |
| guard só `GET` + `POST /admin/x` | passa por fora do guard, sem aviso | aviso | método é parte do match |
| guard descobre o usuário → handler seguinte | impossível | `Locals` | cada rota recebe um `BindRoute` novo, sem estado compartilhado |
| middleware que chama `Write` sem `SetStatus` | `net/http` manda 200 implícito, a cadeia continua e o próximo `SetStatus` vira "superfluous WriteHeader" | ou responde, ou falha | `routeio.Tracked` só marca `answered` em `SetStatus` |
| `HEAD /health` (rota GET) | 405 | 200 sem body | nenhum fallback HEAD→GET |
| guard devolvendo 401/403 | não há constante | `api.StatusUnauthorized` / `StatusForbidden` | `api/server.go` só tem 200/201/204/400/404/405/409/413/415/500 |
| dois `Set-Cookie` | o segundo sobrescreve | acumular | `Response.SetHeader` usa `Header().Set` |
| `?page=2` como `number` | `Entries.Page float64` | `int` | não existe tipo inteiro |
| capturar `/users/{id}` | `add-route` + `add-path id --start 1 --end 1` | uma flag | só há índices de segmento |
| `Ctrl+C` no servidor | processo morre | shutdown gracioso | `Serverdeps.Shutdown` existe e ninguém chama |

## 1. Matching (triggers)

### 1.1 `prefix` passa a respeitar segmentos

| `trigger.type` | Casa `value` = `/admin` com |
|---|---|
| `prefix` (novo) | `/admin`, `/admin/…` — nunca `/administrator` |
| `text-prefix` (novo nome do comportamento atual) | qualquer texto que comece com `/admin` |
| `equal`, `suffix`, `regex` | sem mudança |

`MatchTrigger` (`assets/sandbox-server/sandbox/internal/generated/server/route/IsActionable.go`):

```go
case api.PrefixTrigger:
	value := sandbox.Deps.Stringsdeps.TrimSuffix(trigger.Value, "/")
	return text == value || value == "" ||
		sandbox.Deps.Stringsdeps.HasPrefix(text, value+"/")
case api.TextPrefixTrigger:
	return sandbox.Deps.Stringsdeps.HasPrefix(text, trigger.Value)
```

O `value == ""` cobre `prefix /`, que normaliza para vazio e casa tudo. Isso preserva o logger
do exemplo `route-chain`.

Sem migração: nenhum projeto roda sobre o agnos ainda.

### 1.2 Apelidos na entrada, nome canônico no yaml

`add-route`, `set-route`, `add-path`, `set-path`, `add-parameter` e `set-parameter` aceitam os
apelidos e gravam sempre o nome canônico:

| Apelido | Grava |
|---|---|
| `starts-with` | `prefix` |
| `ends-with` | `suffix` |
| `exact`, `equals` | `equal` |
| `matches` | `regex` |

A normalização fica num único lugar: a função que já produz `unknown trigger type` em
`routeconf`.

### 1.3 Chaves novas no trigger

```yaml
trigger: { type: prefix, value: /admin, negate: false, ignore-case: false }
```

| Chave | Efeito | Flag |
|---|---|---|
| `negate` | inverte o resultado, por exemplo "tudo em `/admin` exceto `/admin/login`" | `--trigger-negate` |
| `ignore-case` | compara em minúsculas (o regex recebe `(?i)`) | `--trigger-ignore-case` |

Para "prefixo X e não Y", usa-se **dois paths** sobre a mesma fatia. Todo path precisa casar, e
isso já é um AND:

```yaml
paths:
  - { id: Route, start: 0, end: -1, trigger: { type: prefix, value: /admin } }
  - { id: NotLogin, start: 0, end: -1, trigger: { type: equal, value: /admin/login, negate: true } }
```

`api.Trigger` ganha `Negate bool` e `IgnoreCase bool`, e o `new.go` gerado passa a escrevê-los.

## 2. Métodos

| Mudança | Onde |
|---|---|
| `methods: [ANY]` casa qualquer método. Na entrada, `*` é apelido de `ANY` | `acceptsMethod`, `accepts` (servermain), `routeconf`, `check_routes` |
| `ANY` não pode ser combinado com outro método | `check_routes` |
| `HEAD` sem rota `HEAD` declarada roda a rota `GET` com um `Write` descartado | `dispatch`: um segundo passe quando `method == HEAD` e nada casou |
| `OPTIONS` continua sem tratamento automático (ver §8, CORS) | — |

## 3. Prioridade e cadeia

### 3.1 Faixas padrão em vez de tudo no 0

| Criado por | `priority` padrão |
|---|---|
| `add-route` | `100` |
| `add-route --middleware` | `10` |
| `add-page` | `100` |
| `health` (template do `server-init`) | `100` |

A regra de ser `≥ 0` continua. Só o valor padrão muda, o que abre
espaço abaixo das rotas para inserir guards sem renumerar nada.

### 3.2 Prioridade relativa

`add-route` e `set-route` ganham:

| Flag | Grava |
|---|---|
| `--before <route>` | `priority` da rota alvo − 1 (erro se isso ficar `< 0`, sugerindo `rebalance-routes`) |
| `--after <route>` | `priority` da rota alvo + 1 |

Essas flags não podem ser combinadas com `--priority`. O yaml continua guardando um inteiro
absoluto, sem referência a outra rota. Assim o collector não muda e não existem ciclos.

### 3.3 `set-route --priority` sem sentinela

Hoje o padrão `-1` significa "não mexer". O certo é tirar o `default` e usar a presença da flag,
como `--help` já faz. Assim `-1` passa a ser rejeitado com a mensagem certa, e não ignorado em
silêncio.

### 3.4 `add-route --middleware`

Só muda o que é escrito. **Não existe chave `kind` no yaml**: um middleware continua sendo uma
rota que não escreve status.

| Chave | Valor gravado |
|---|---|
| `methods` | `[ANY]` |
| `priority` | `10` |
| `trigger-type` padrão | `prefix` |
| `response-type` | continua obrigatório (a regra não muda); grava `text/plain` |
| stub `InternalPureHandler.go` | usa `Locals` e retorna `nil` sem status (ver §4) |

### 3.5 "Escrever é responder" — corrige o 200 implícito

Em `routeio.Tracked`, um `Write` sem status antes chama `SetStatus(api.StatusOk)`. Com isso,
`answered` passa a ser verdadeiro. A regra da cadeia vira: **uma rota responde quando define um
status ou escreve bytes.** `SetHeader` sozinho continua não respondendo, e é assim que um
middleware adiciona headers.

Precisa ser atualizada nos docs: `RouteYaml#the-chain`, `Rules`, `ServerUsage` e no bullet
"The server is a chain" do `CLAUDE.md`.

### 3.6 Fase `after` (opcional, depois do resto)

```yaml
phase: after   # default: before
```

As rotas `after` que casam rodam **depois** que a cadeia respondeu, em ordem de `priority`, e
nunca respondem (um `SetStatus` ali é ignorado e gera log). Elas cobrem log de status e latência
e métricas. Para ler o status que foi enviado, `routeio.Tracked` passa a guardar o `status` além
de `answered`. Flag: `add-route --phase after`.

## 4. Passagem de dados na cadeia — `Locals`

```go
// api.Route
// Locals is one request's scratch space, shared by every route of the chain
// that runs for it: what a middleware stores, the routes after it read.
Locals map[string]any
```

- `dispatch` cria `locals := map[string]any{}` **uma vez por request** e faz
  `bound.Locals = locals` em cada `BindRoute`.
- Handlers de erro (`failRequest`, `recoverRoute`) recebem o mesmo mapa.
- Helpers tipados em `assets/sandbox-server/sandbox/internal/generated/routeio/locals.go`:

```go
func SetLocal(route *api.Route, key string, value any)
func GetLocal[T any](route *api.Route, key string) (T, bool)
```

Exemplo, guard + rota:

```go
// routeslist/admin_guard/InternalPureHandler.go — priority 10, methods [ANY], prefix /admin
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	user, ok := auth.Check(sandbox, entries.Authorization)
	if !ok {
		return routeio.Fail(sandbox, route, api.StatusUnauthorized, "authorization", "invalid token")
	}
	routeio.SetLocal(route, "user", user)
	return nil
}
```

Para `Fail` com 401/403 chegar a um arquivo do projeto, veja §6.2.

## 5. Declaração de caminhos e argumentos

### 5.1 `--pattern` em `add-route`

```bash
agnos add-route get-article --pattern '/get-article/{article}'
agnos add-route get-post    --pattern '/users/{user:uuid}/posts/{post:integer}'
agnos add-route admin-area  --pattern '/admin/{*rest}'
```

`--pattern` é convertido para as chaves que já existem no yaml, e o yaml não ganha uma chave
`pattern`. O `self.Pattern` gerado continua existindo só para os docs. `--pattern` não pode ser
combinado com `--trigger` e `--trigger-type`.

| Trecho do pattern | Vira |
|---|---|
| literais consecutivos, por exemplo `/get-article` nos índices i..j | `{id: Seg<i>, start: i, end: j, trigger: {type: equal, value: /get-article}}` |
| `{article}` | `{id: Article, start: i, end: i}` → `Entries.Article string` |
| `{article:<tipo>}` | o mesmo path com `type: <tipo>` → `Entries.Article` no tipo Go da tabela §5.2 |
| `{*rest}` (só no fim) | `{id: Rest, start: i, end: -1}` → `Entries.Rest string`, por exemplo `/a/b.png` |
| sem `{*…}` no fim | `segments: <n>` na rota: a request precisa ter exatamente n segmentos |
| com `{*…}` no fim | `segments` omitido, e o `{*rest}` exige pelo menos um segmento |

`segments` é uma chave nova da rota (int, opcional, ≥ 1), editada por `set-route --segments` e
removida com `--clear segments`. Sem ela, `/get-article/42/extra` também casaria, porque um path
só olha a sua própria fatia.

Exemplo completo, `add-route get-article --pattern '/get-article/{article:integer}'`:

```yaml
methods: [GET]
priority: 100
response-type: application/json
segments: 2
paths:
  - { id: Seg0, start: 0, end: 0, trigger: { type: equal, value: /get-article } }
  - { id: Article, start: 1, end: 1, type: integer }
```

```go
type Entries struct {
	FullRoute string `id:"FullRoute"`
	Seg0      string `id:"Seg0"`
	Article   int    `id:"Article"`
}
```

```
GET /get-article/42      Entries.Article = 42
GET /get-article/abc     404, a rota não roda (o tipo faz parte do match)
GET /get-article/42/x    404, segments: 2
GET /get-article         404, sem segmento 1
```

### 5.2 Tipo em path

`paths[].type` define o tipo Go do campo em `Entries` e faz parte do match. Um segmento que não
converte é **não-match** (a cadeia segue e, sem ninguém responder, 404), e não um 400: pela
mesma regra do `trigger`, a URL não é desta rota. Assim, `/articles/{id:integer}` e
`/articles/{slug}` convivem, com o integer numa `priority` menor.

| `type` | Go | Aceita |
|---|---|---|
| `string` (padrão) | `string` | qualquer segmento |
| `integer` | `int` | `^-?\d+$` |
| `number` | `float64` | o que `ParseFloat` aceita |
| `uuid` | `string` | UUID canônico, 8-4-4-4-12 hex |

O tipo só vale para um path de um único segmento (`start == end`). `verify` recusa um `type`
diferente de `string` num path de vários segmentos ou com `end: -1`. Para uma rota que já existe,
as flags são `add-path --type` e `set-path --type`.

A conversão fica em `IsActionable` (decide o match) e `RequestHandler` (preenche `Entries`), as
duas usando a mesma função `parsePathValue`, para as duas nunca divergirem.

### 5.3 Tipos de parâmetro

| Novo `type` | Go | Motivo |
|---|---|---|
| `integer` | `int` | `page`, `limit` e `offset` hoje saem como `float64` |
| `integer-array` | `[]int` | par de `string-array` |

### 5.4 Fonte `cookie`

`fonts: [cookie]` lê `Cookie: <key>=…`. Precisa de `Request.GetCookie func(name string) string`
em `serverdeps` (contrato e adapter em `assets/deplist/` e `assets/adapterlist/`, com o
byte-for-byte da cópia deste repo).

### 5.5 Body

| `body.type` novo | `ReadBody` retorna |
|---|---|
| `form` (`application/x-www-form-urlencoded`) | `(map[string][]string, error)` |
| `multipart` | fica para depois; exige contrato de arquivo em `serverdeps` |

## 6. Resposta e contratos

### 6.1 `serverdeps`

| Membro novo | Motivo |
|---|---|
| `Response.AddHeader(key, value)` | múltiplos `Set-Cookie` e `Vary` |
| `Response.GetHeader(key)` | um middleware `after` ou uma rota lê o que outro já definiu |
| `Request.GetCookie(name)` | §5.4 |
| `Request.GetHost()` | redirects absolutos, multi-host |
| `Request.GetHeaders() map[string][]string` | logging e proxy |

Os deps continuam sem opinião: são só acesso ao `net/http`, sem política.

### 6.2 Status e seus `Handle*`

Constantes novas em `api/server.go`: `StatusMovedPermanently 301`, `StatusFound 302`,
`StatusSeeOther 303`, `StatusNotModified 304`, `StatusTemporaryRedirect 307`,
`StatusPermanentRedirect 308`, `StatusUnauthorized 401`, `StatusForbidden 403`,
`StatusUnprocessable 422`, `StatusTooManyRequests 429`, `StatusUnavailable 503`.

`server.Fail` ganha dois ramos e dois arquivos de `sandbox/internal/server/errors/`, escritos uma
vez como os outros seis:

| Status | Arquivo |
|---|---|
| 401 | `handle_unauthorized.go` |
| 403 | `handle_forbidden.go` |

`verify` passa a checar oito arquivos.

### 6.3 Helpers de resposta em `routeio`

```go
func WriteJSON(sandbox *api.Sandbox, response serverdeps.Response, status int, value any) error
func WriteText(response serverdeps.Response, status int, text string) error
func Redirect(response serverdeps.Response, status int, location string) error
```

Hoje cada handler faz à mão serialização, `SetStatus` e `Write`.

## 7. Comandos

### 7.1 Novos

| Comando | Faz | Escreve |
|---|---|---|
| `list-routes` | tabela da cadeia **na ordem de execução**: priority, nome, methods, pattern, fase | nada |
| `explain-route <METHOD> <path> [--header k=v]… [--query k=v]…` | roda `IsActionable` estaticamente sobre os `route.yaml` e imprime, por rota, `runs` ou `skipped: <path X trigger falhou / método>`, mais o desfecho (404/405) quando nada casa | nada |
| `rename-route <old> <new>` | move o diretório, reescreve o `package` do `InternalPureHandler.go` e roda `build` | sim |
| `rebalance-routes [--step 10]` | redistribui as prioridades mantendo a ordem atual | todos os `route.yaml` |

`explain-route` é o comando que responde "por que meu `/admin` não pegou?" sem subir o servidor.
Ele reutiliza `MatchTrigger`/`PathSlice`. Para isso, essas funções precisam existir também do
lado do agnos (`sandbox/internal/actions/explain_route/`), já que o template de `IsActionable` é
código do projeto gerado e não do agnos. Isso é duplicação consciente: um golden de exemplo
precisa travar os dois lados iguais (ver §10).

### 7.2 Alterados

| Comando | Mudança |
|---|---|
| `add-route` | `--help`/`--category` deixam de ser obrigatórios (padrões: `""` e `Routes`); `--pattern`; `--middleware`; `--before`/`--after`; `--phase`; `--trigger-negate`, `--trigger-ignore-case`; padrão de `priority` `100`; aceita `ANY` e apelidos de trigger |
| `set-route` | `--priority` sem sentinela; `--before`/`--after`; `--segments`; `--phase`; `--clear-segments` |
| `add-path` / `set-path` | `--type`; `--trigger-negate`; `--trigger-ignore-case`; apelidos |
| `add-parameter` / `set-parameter` | `--type integer\|integer-array`; `--font cookie`; `--trigger-negate`; `--trigger-ignore-case`; apelidos |
| `set-body` | `--type form` |
| `show-route` | mostra também a posição na cadeia ("#3 de 7, depois de `admin-guard`") |
| `add-page` | padrão de `priority` `100` |

## 8. Processo do servidor

| Mudança | Onde |
|---|---|
| `SIGINT`/`SIGTERM` chamam `Serverdeps.Shutdown` e a request em andamento termina | contrato novo `Signaldeps.OnInterrupt(func())` (dep + adapter + available), chamado em `ServerMain` |
| `start-server --shutdown-timeout-ms` | `ServeProps.ShutdownTimeoutMs` |
| middleware CORS opcional | `server-init --cors` escreve `routeslist/cors/` (priority 0, `ANY`, prefix `/`) como uma rota comum, que o projeto passa a dono. Não é mecânica do dispatch |

## 9. Regras do `verify` (`check_routes.go`)

| Regra nova |
|---|
| `ANY` sozinho em `methods` |
| `trigger.type` ∈ {`equal`, `prefix`, `text-prefix`, `suffix`, `regex`} (o apelido nunca chega ao yaml) |
| `negate`/`ignore-case` só booleanos |
| `paths[].type` ∈ {`string`, `integer`, `number`, `uuid`}, e diferente de `string` só com `start == end` |
| `segments` ≥ 1 quando declarado |
| `phase` ∈ {`before`, `after`}; rota `after` sem `body` |
| existem os 8 `handle_*.go` |
| **aviso** (não erro): rota com prefix/`ANY` e `priority` menor que outra, cujos métodos não cobrem todos os métodos das rotas que ela alcança (o caso "guard só GET") |

## 10. Exemplos (`examples/`)

| Novo exemplo | Trava |
|---|---|
| `route-prefix-segment` | `/admin` × `/administrator` |
| `route-trigger-alias` | `starts-with` gravado como `prefix` |
| `route-middleware-guard` | `--middleware`, `Locals`, 401 |
| `route-before-after` | `--before`/`--after` e `rebalance-routes` |
| `route-pattern` | `--pattern` com `{x}`, `{x:integer}`, `{x:uuid}` e `{*rest}`; `/get-article/abc` → 404; `segments` |
| `explain-route` | saída para match, skip por trigger, 405 e 404 |
| `route-head-fallback` | HEAD→GET |

Os exemplos `route-chain`, `server-route` e `server-init` mudam de golden por causa do padrão
`100` e dos dois `handle_*.go` novos. Rodar `exec-test --update` só depois de conferir o diff.

## 11. Docs

| Página (template) | Mudança |
|---|---|
| `assets/doc-server/docs/RouteYaml/doc.md` | triggers, `negate`/`ignore-case`, `ANY`, `segments`, `phase`, `Locals`, "escrever é responder", tipos novos, os 8 `Handle*` |
| `assets/doc-server/docs/ServerUsage/doc.md` | receita de guard e de `--pattern` |
| `assets/doc/docs/Rules/doc.md` | regra da cadeia atualizada; faixas de priority |
| `assets/doc/docs/Workflow/doc.md` | `explain-route` como passo de depuração |
| `docs/Contributing/doc.md` | padrão de middleware |
| `CLAUDE.md` | bullet "The server is a chain" e a lista dos `handle_*.go` (6 → 8) |

## 12. Ordem de implementação

| # | Pacote | Por quê primeiro |
|---|---|---|
| 1 | §1.1, §1.2 (prefix por segmento + apelidos) | bug real, mudança pequena |
| 2 | §3.5 (escrever é responder) | bug real, afeta a semântica de toda cadeia |
| 3 | §2 `ANY` + §3.1 faixas + §3.4 `--middleware` + §6.2 401/403 | torna um `/admin` protegido possível sem gambiarra |
| 4 | §4 `Locals` | torna o middleware útil |
| 5 | §7.1 `list-routes` + `explain-route` | depuração da cadeia |
| 6 | §5.1–5.3 (`--pattern`, `segments`, `integer`) | ergonomia de declaração |
| 7 | §3.2 `--before`/`--after`, `rebalance-routes`, `rename-route` | ergonomia de prioridade |
| 8 | §2 HEAD, §6.1 contratos, §6.3 helpers, §5.4 cookie | completude |
| 9 | §1.3 `negate`/`ignore-case`, §3.6 `after`, §5.5 `form`, §8 | o resto |

Cada pacote fecha com: `bootstrap build`, `build -q && git diff --quiet` (idempotência),
`exec-test --only <exemplos do pacote>`, e o golden do exemplo novo.
