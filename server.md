# Plano de implementação — mecânica de servidores

Camada `server`, espelho exato da camada `cli`. Cada arquivo abaixo é a instância *server* de
um arquivo *cli* que já existe (regra "Every file is an instance of a pattern").

## 1. Correspondência com a camada CLI

| Conceito | CLI (existe) | Server (a criar) |
|---|---|---|
| Contrato de entrada externa | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` |
| Adapter | `adapters/libs/verb/` | `adapters/libs/serverdeps/` |
| Superfície pública | `sandbox/api/cli.go` | `sandbox/api/server.go` |
| Bind | `sandbox/binds/cli.go` | `sandbox/binds/server.go` |
| Despacho gerado | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/servermain.go` |
| Unidade declarada | `sandbox/internal/commands/<name>/` | `sandbox/internal/routes/<name>/` |
| Declaração (yaml) | `entries.yaml` | `route.yaml` |
| Struct gerada | `entries.go` -> `Entries` | `entries.go` -> `Entries` |
| Handler (à mão) | `handler.go` -> `CommandHandler` | `handler.go` -> `RouteHandler` |
| Parsable | `parsables/commandconf/` | `parsables/routeconf/` |
| Grupo de assets | `assets/cli/` | `assets/server/` |
| Instalação / remoção | `cli-init` / `cli-purge` | `server-init` / `server-purge` |
| Editores da declaração | `add-command`, `add-flag`, `add-arg`, `set-command` | `add-route`, `set-route`, `add-segment`, `add-header`, `add-param`, `set-body`, `add-body-field` (+ inversos) |
| Gatilho no build | `hasCli := io.IsDir("sandbox/internal/cli")` | `hasServer := io.IsDir("sandbox/internal/server")` |

Invariantes:

- **O despacho valida tudo que não é o corpo antes do handler rodar.** Um `RouteHandler` nunca
  *decide* `400`/`404`/`405`/`413`/`415` — quem responde isso é `servermain.go`.
- **O corpo é a exceção:** não é lido nem parseado antes do handler. `Entries` não tem campo
  `Body`; tem o método `ReadBody`, que lê, valida e converte sob demanda, e já responde o erro
  devolvendo o status para o handler propagar.

---

## 2. Fase 1 — dep `serverdeps`

Campo `Deps.Serverdeps`, lib adapter `serverdeps`, sobre `net/http` (stdlib -> **sem** entrada
em `assets/depsversion.yaml`). O dep é **0 opinativo**: abre a porta, aplica timeouts e entrega
toda requisição à única `ServerProps.Handler`. Rota, método, `{param}`, `404` e `405` são do
`servermain.go` gerado — sem `ServeMux`, sem `PathValue`.

| Arquivo | Papel |
|---|---|
| `assets/deplist/serverdeps/sandbox/deps/serverdeps/serverdeps.go` | template do contrato |
| `assets/deplist/serverdeps/adapters/libs/serverdeps/serverdeps.go` | template do adapter |
| `sandbox/deps/serverdeps/serverdeps.go` | cópia instalada (byte-a-byte com o template — **(verify)** `check_deplist.go`) |
| `adapters/libs/serverdeps/serverdeps.go` | idem |

```go
package serverdeps

type Lib struct {
	NewServer func(props ServerProps) Server
}

type ServerProps struct {
	Addr           string // ":8080"
	ReadTimeoutMs  int
	WriteTimeoutMs int
	Handler        func(Request, Response) // uma vez por requisição, qualquer que seja
}

type Server struct {
	Listen   func() error // bloqueia até Shutdown ou erro
	Shutdown func() error
}

type Request struct {
	GetMethod     func() string
	GetPath       func() string // caminho cru; quem o fatia em segmentos é o sandbox
	GetHeader     func(key string) string
	GetQueryParam func(name string) string
	GetQueryAll   func(name string) []string // campos `array: true`
	ReadBody      func(limit int) ([]byte, error) // limit -1 = tudo
	GetRemoteAddr func() string
}

type Response struct {
	SetHeader func(key string, value string)
	SetStatus func(code int)
	Write     func(body []byte) error
}
```

Só builtins atravessam o contrato — nada de `time.Time`, `io.Reader` ou tipos de `net/http`
(mesma razão do commit `71cc554`). Toda declaração exportada leva doc comment.
Editar também `assets/all/docs/DepList/doc.md` (uma linha).

---

## 3. Fase 2 — `route.yaml`, a declaração

`sandbox/internal/routes/<name>/route.yaml`, escrito por `add-route` e reescrito por um editor
por lugar que o arquivo guarda algo — `set-route`, `add-segment`, `add-header`, `add-param`,
`set-body`, `add-body-field` e seus inversos — **nunca à mão**. Diretório snake_case para rota
kebab-case (`get-user` -> `get_user/`).

Não existe chave `path`: o caminho é a concatenação dos segmentos de `paths`, derivada pelo
collector. Um segmento é de um de dois tipos:

| Segmento | Chave | Papel |
|---|---|---|
| **trigger** | `identifier` | literal na URL, **sempre começando com `/`** (`/users`, nunca `users`). O primeiro `identifier` da rota é o gatilho que a nomeia em `docs/Routes` |
| **captura** | `name` | segmento variável, vira campo de `Entries` já convertido |

**Todo `identifier` começa com `/`.** A barra inicial é obrigatória: é ela que faz o segmento se
ler como caminho no arquivo, nas mensagens e nos docs, sem o leitor reconstruir mentalmente a
concatenação. `/` sozinho é o identifier da raiz; nenhum outro tem `/` interno ou final, e
nenhum é vazio **(verify)**.

**Ordem de casamento — da rota mais específica para a menos.** A ordem em que o despacho testa
as rotas não é a ordem do diretório: o collector ordena por

1. **mais `identifier`s primeiro** — a rota que fixa mais segmentos literais é a mais específica;
2. empate: **maior soma de caracteres dos `identifier`s**;
3. empate: `Pattern()` em ordem alfabética (determinismo).

Sem isso uma rota `/` e uma rota `/home` colidem: avaliada primeiro, `/` casa e `/home` nunca é
alcançada. A ordenação é do collector, não do template — `.Routes` já chega ordenado a
`servermain.go` **(verify não precisa checar, o build garante)**.

O mesmo `name` pode vir de mais de uma origem — o campo de `Entries` é único e vence a primeira
origem declarada que trouxer valor (`paths`, `headers`, `params`, na ordem do arquivo).
`required` é satisfeito por qualquer uma; o tipo tem de ser o mesmo em todas **(verify)**.

```yaml
# sandbox/internal/routes/create_user/route.yaml -> POST /users/{tenant}/create
method: POST
paths:
  - identifier: "/users"
  - name: "tenant"
    type: string
    required: true
  - identifier: "/create"
category: Users
help: Create a user under a tenant
examples:
  - "curl -X POST localhost:8080/users/acme/create -d '{\"email\":\"a@b.c\",\"age\":30}'"
headers:
  - name: authorization
    type: string
    required: true
params:
  - name: page
    type: int
    default: "1"
    min: 1
    max: 100
  - name: tag
    type: string
    array: true
body:
  type: json            # none (padrão) | raw | text | json
  required: true
  max-bytes: 1048576
  content-type: application/json
  json-schema:
    type: object
    additionalProperties: false
    required: [email, age]
    properties:
      email: { type: string, format: email, minLength: 3, maxLength: 254 }
      age:   { type: integer, minimum: 0, maximum: 130 }
      tags:  { type: array, items: { type: string } }
      address:
        type: object
        required: [city]
        properties:
          city: { type: string }
          zip:  { type: string, pattern: "^[0-9]{5}$" }
```

### Chaves

| Chave de rota | Efeito |
|---|---|
| `paths` | segmentos da URL, na ordem. Obrigatória e não vazia |
| `method` | `GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS`. Default `GET` |
| `category`, `help`, `long-description`, `examples`, `hidden` | idêntico a `entries.yaml`; alimenta `docs/Routes` |
| `headers`, `params` | sequências de campos lidos do cabeçalho / da query |
| `body` | descrição do corpo (objeto, não sequência) |

| Chave de segmento | Efeito |
|---|---|
| `identifier` | literal começando com `/` (`/users`). Exclui `name` e toda chave de campo no mesmo item |
| `name` | segmento capturado; posição variável no casamento e campo de `Entries` |
| `type` | `string`\|`boolean`\|`int`\|`float`. Default `string` |
| `description`, `examples`, `min`, `max` | como em `entries.yaml` |
| `required` | sempre `true` — `false` é erro de **(verify)**. `array` e `default` proibidos |

**Campos de `headers`/`params`:** mesmas chaves de `entries.yaml` (`name`, `description`,
`examples`, `type`, `default`, `required`, `array`, `min`, `max`). O `name` **é** a grafia
externa: nome do header (casado sem diferenciar maiúsculas) ou chave na query. `array: true` só
em `params`.

| Chave de `body` | Efeito |
|---|---|
| `type` | `none` (padrão), `raw` (`[]byte`), `text` (`string`), `json` |
| `required` | corpo ausente/vazio vira `400` |
| `max-bytes` | excedido vira `413`. Default 1 MiB |
| `content-type` | divergência vira `415`. Default `application/json` para `type: json` |
| `json-schema` | subconjunto de JSON Schema, só quando `type: json` |

Do despacho, só `content-type` (é header) e um `Content-Length` já maior que `max-bytes`; o
resto é aplicado dentro de `ReadBody`.

**Subconjunto de JSON Schema:** `type` (`object`/`array`/`string`/`integer`/`number`/`boolean`/
`null`), `properties`, `required`, `additionalProperties` (bool), `items`, `enum`, `const`,
`minimum`, `maximum`, `exclusiveMinimum`, `exclusiveMaximum`, `minLength`, `maxLength`,
`pattern`, `minItems`, `maxItems`, `uniqueItems`, `format` (`email`, `uuid`, `date-time`,
`uri`), `nullable`. Fora dele (`$ref`, `oneOf`, `allOf`, `anyOf`, `patternProperties`) -> erro
de `verify`/`build`, não silêncio.

---

## 4. Fase 3 — parsable `routeconf`

`sandbox/internal/parsables/routeconf/` — mesmos nomes e ordem que `commandconf/`:

| Arquivo | Conteúdo |
|---|---|
| `api.go` | `RouteConf` (com `Paths []Segment`), `Segment`, `Field`, `Body`, `Schema` |
| `new.go` | `New(deps, content) (*RouteConf, error)` |
| `new_empty.go` | `NewEmpty(deps) *RouteConf` |
| `bind_methods.go` | `BindMethods` (liga `Render`) |
| `render.go` | `Render` -> `route.yaml` canônico, chaves em ordem fixa (idempotência) |

`Segment` é `{Identifier string; Field *Field}`, exatamente um preenchido. `RouteConf.Pattern()`
é a concatenação dos segmentos na ordem — o `identifier` já traz a própria `/`, a captura entra
como `/{name}` — usado em `docs/Routes` e nas mensagens; o casamento compara segmento a
segmento, não o texto. `RouteConf.IdentifierCount()` e `RouteConf.IdentifierLen()` alimentam a
ordenação de casamento. `Schema` é árvore recursiva (`Type`, `Properties []SchemaProperty`,
`Items *Schema`, `Required []string`, bounds com o par valor/`Has…` como em
`Field.Min`/`HasMin`).

---

## 5. Fase 4 — a camada server (grupo `assets/server/`)

Renderizado por todo `build` quando `sandbox/internal/server/` existe; todo arquivo do grupo é
sempre reescrito.

| Template (`assets/server/`) | Destino |
|---|---|
| `sandbox/api/server.go` | `sandbox/api/server.go` |
| `sandbox/binds/server.go` | `sandbox/binds/server.go` |
| `sandbox/internal/server/servermain.go` | idem |
| `sandbox/internal/routeio/{jsonschema,write_error}.go` | idem |
| `sandbox/internal/routes/health/{route.yaml,handler.go}` | idem — rota embutida, análoga ao comando `version` |
| `docs/{RouteYaml,Routes,ServerUsage}/{doc.md,props.yaml}` | `docs/…` (`Routes` renderizado das rotas, como `docs/Commands`) |

### `sandbox/api/server.go`

```go
type Server struct {
	Serve func(props ServeProps) error
}

type ServeProps struct {
	Addr           string
	ReadTimeoutMs  int
	WriteTimeoutMs int
}

const (
	StatusOk               = 200
	StatusCreated          = 201
	StatusNoContent        = 204
	StatusBadRequest       = 400
	StatusNotFound         = 404
	StatusMethodNotAllowed = 405
	StatusConflict         = 409
	StatusPayloadTooLarge  = 413
	StatusUnsupportedMedia = 415
	StatusFailure          = 500
)
```

`sandbox/api/sandbox.go` (campo `Server`) e `sandbox/new.go` (chamada a `binds.ServerBind`) já
são gerados por convenção — nada a editar à mão.

### `servermain.go` (gerado)

Uma única função entregue ao dep; `ServerMain` não tem `range`, o `{{range .Routes}}` está no
`dispatch`:

```go
func ServerMain(deps *deps.Deps, props ServeProps) error {
	server := deps.Serverdeps.NewServer(serverdeps.ServerProps{
		Addr:           props.Addr,
		ReadTimeoutMs:  props.ReadTimeoutMs,
		WriteTimeoutMs: props.WriteTimeoutMs,
		Handler: func(req serverdeps.Request, res serverdeps.Response) {
			dispatch(deps, req, res)
		},
	})
	return server.Listen()
}

func dispatch(deps *deps.Deps, req serverdeps.Request, res serverdeps.Response) {
	segments := splitPath(deps, req.GetPath())
	method := req.GetMethod()
	path_matched := false
{{- /* .Routes chega ordenado: mais identifiers primeiro, depois identifiers mais longos */}}
{{- range .Routes}}
	if match{{.GoName}}(segments) {
		path_matched = true
		if method == "{{.Method}}" {
			handle{{.GoName}}(deps, req, res, segments)
			return
		}
	}
{{- end}}
	if path_matched {
		routeio.WriteError(deps, res, api.StatusMethodNotAllowed, "", "method not allowed")
		return
	}
	routeio.WriteError(deps, res, api.StatusNotFound, "", "route not found")
}
```

Mais, por rota: `match<GoName>(segments []string) bool` (comprimento + `identifier` por
posição — comparado sem a `/` inicial, que o `splitPath` já consumiu — capturas aceitando
qualquer valor) e `handle<GoName>`, que converte segmentos capturados (por índice), headers e
params na ordem de declaração; aplica defaults; checa
`required`, `min`/`max`, `array`, `content-type` e `Content-Length`; monta `Entries` (com
`Entries.Request`); e chama `routes_<name>.RouteHandler(deps, &entries, res)`. **Nada é lido do
socket aqui.** Falhas saem por `routeio.WriteError` em JSON (`{"error": "...", "field": "..."}`)
e registram em `deps.Std.Log`. Helpers (`readHeader`, `parseIntValue`, `bindQueryArray`,
`splitPath` com `deps.Stringsdeps.Split`) ficam no rodapé, como os de argv em `climain.go`.

| Situação | Status | Quem responde |
|---|---|---|
| nenhum `match<GoName>` casou | 404 | despacho |
| path casou, método divergente | 405 | despacho |
| `content-type` divergente | 415 | despacho |
| `Content-Length` acima de `max-bytes` | 413 | despacho |
| header/param inválido, `required` ausente, fora de `min`/`max` | 400 | despacho |
| corpo excede `max-bytes` na leitura | 413 | `ReadBody` |
| corpo ausente com `required: true`, JSON inválido, schema reprovado | 400 | `ReadBody` |
| handler devolveu `0` ou entrou em pânico | 500 | despacho |

### `routeio/` (gerado, estático)

Pacote próprio porque `servermain.go` importa todas as rotas — uma rota não pode importar
`internal/server` de volta. `routeio` não importa nenhum dos dois e é importado pelos dois.

```go
// jsonschema.go — validador puro do subconjunto, sobre deps.Serializables.ParseJson;
// devolve o objeto parseado, a mensagem da primeira violação, e se passou.
func ValidateSchema(deps *deps.Deps, schema_json string, body []byte) (*serializables.SerializibleObject, string, bool)

// write_error.go — única forma de escrever uma falha, usada pelo despacho e por todo ReadBody;
// devolve o próprio status, para o chamador retornar em uma linha.
func WriteError(deps *deps.Deps, response serverdeps.Response, status int, field string, message string) int
```

---

## 6. Fase 5 — geração por rota, no `build`

| Arquivo | Papel |
|---|---|
| `actions/build/collect_routes.go` | `CollectRoutes(deps, io)`: lê cada `route.yaml` via `routeconf` -> `[]map[string]any` (`Name`, `GoName`, `Method`, `Trigger`, `Path`, `Segments`, `Headers`, `Params`, `Body`, `SchemaJson`, `BodyStructs`, `HasBody`). Cópia de `collect_commands.go`. Devolve a lista **ordenada para o casamento** (mais `identifier`s, depois soma de caracteres dos `identifier`s, depois `Pattern()`), via `deps.Sortdeps` |
| `actions/build/collect_route_docs.go` | alimenta `docs/Routes` — cópia de `collect_command_docs.go` |
| `actions/build/generate_route_entries.go` | renderiza `assets/templates/route_entries.go` por rota em `routes/<name>/entries.go` |
| `assets/templates/route_{entries.go,route.yaml,handler.go}` | template do `Entries` e scaffolds usados por `add-route` |
| `assets/templates/start_server_{entries.yaml,handler.go}` | scaffold do comando `start-server`, escrito por `server-init` |

Em `actions/build/build_internal.go`: `hasServer := io.IsDir("sandbox/internal/server")`,
chamadas a `CollectRoutes`/`CollectRouteDocs`, vars `"HasServer"`, `"Routes"`, `"RouteDocs"`, e
ao final `if hasServer { GenerateRouteEntries(...); utils.RenderGroup(deps, io, "server", vars) }`.
`GeneratedDocsGroups(hasCli, hasServer)` em `collect_generated_docs.go` passa a devolver
`"server"` também.

### `entries.go` gerado

```go
package create_user

type Entries struct {
	Tenant        string             // segmentos capturados, na ordem declarada
	Authorization string             // headers, na ordem declarada
	Page          int                // params, na ordem declarada
	Tag           []string
	Request       serverdeps.Request // preenchido pelo despacho; fonte de ReadBody
}

type Body struct {
	Email   string
	Age     int
	Tags    []string
	Address BodyAddress
}

type BodyAddress struct {
	City string
	Zip  string
}

const EntriesSchema = `{"type":"object", ...}`   // schema canônico, para ValidateSchema

// ReadBody lê, valida e converte o corpo na primeira chamada (cache nas seguintes). Devolve
// api.StatusOk quando passou; nos demais status a resposta de erro já foi escrita.
func (entries *Entries) ReadBody(deps *deps.Deps, response serverdeps.Response) (Body, int)
```

| `body.type` | Assinatura gerada |
|---|---|
| `none` (padrão) | nenhum `ReadBody` |
| `raw` | `([]byte, int)` |
| `text` | `(string, int)` |
| `json` com `json-schema` de objeto | `(Body, int)` |
| `json` sem `json-schema` | `(*serializables.SerializibleObject, int)` |

Toda variante faz, na ordem: `Request.ReadBody(max-bytes)` (`413`), `required` (`400`) e — só em
`json` — `routeio.ValidateSchema` contra `EntriesSchema` (`400` na primeira violação, com o
campo em `"field"`). Segunda chamada devolve o valor em cache.
Nomes: objeto aninhado -> `Body<Caminho>`; item de array de objetos -> `Body<Caminho>Item`.

### `handler.go` (único arquivo de uma rota escrito à mão)

```go
func RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int {
	// checagens do handler primeiro — o corpo ainda não foi lido
	if !isAuthorized(deps, entries.Authorization) {
		return writeJson(response, api.StatusFailure, ...)
	}
	body, status := entries.ReadBody(deps, response) // ausente se body.type: none
	if status != api.StatusOk {
		return status
	}
	return writeJson(response, api.StatusCreated, createUser(deps, entries.Tenant, body))
}
```

Retorna o status com que respondeu (análogo ao `int` de `CommandHandler`); o único caminho pelo
qual devolve `400`/`413`/`415` é propagando o de `ReadBody`.

---

## 7. Fase 6 — ações e comandos

Duas camadas por feature: `<name>.go` (abre SmartIO, persiste, dispara `build`) +
`<name>_internal.go` (lógica sobre SmartIO aberto). Cada `commands/<dir>/` traz `entries.yaml`
(via `add-command` + `add-flag`/`add-arg`), `entries.go` (gerado) e `handler.go` (à mão).

| Ação / comando | Faz |
|---|---|
| `server_init` -> `server-init` | roda `cli-init` quando o projeto não tem CLI; instala `serverdeps`, `std`, `serializables`, `stringsdeps`; renderiza o grupo `server`; escreve o comando `start-server`; roda `build` |
| `server_purge` -> `server-purge` | remove o grupo `server` + `sandbox/internal/{server,routes}` (cópia de `cli_purge`) |
| `add_route` -> `add-route` | escreve `route.yaml` + `handler.go` (recusa sobrescrever) |
| `remove_route` -> `remove-route` | apaga o diretório da rota |
| `set_route` -> `set-route` | reescreve chaves de nível de rota (`method`, `help`, `category`, `long-description`, `hidden`, `examples`) |
| `add_segment` / `remove_segment` | acrescenta ou remove um segmento de `paths` |
| `add_header` / `remove_header` | declara ou remove um header de requisição |
| `add_param` / `remove_param` | declara ou remove um parâmetro de query |
| `set_body` -> `set-body` | reescreve o envelope do corpo (`type`, `required`, `max-bytes`, `content-type`, `--drop-schema`) |
| `add_body_field` / `remove_body_field` | declara ou remove uma propriedade do `json-schema` |

Um editor por origem, espelhando `add-flag`/`add-arg` do lado cli: nenhuma chave da declaração
precisa de edição à mão. `add-segment` acrescenta ao fim de `paths` (ou em `--position`), com
`--identifier` para trigger ou um nome para captura; `--identifier` normaliza o valor para
começar com `/` (`users` -> `/users`) e recusa `/` interno ou final. `add-body-field` aceita
caminho pontuado (`address.city`), criando os objetos intermediários no `json-schema`, e cobre
todo o subconjunto de keywords (`enum`, `const`, `nullable`, `additionalProperties`, os limites
exclusivos e os de lista).

### `server-init` implica CLI

Um servidor precisa de um ponto de entrada que o suba, e esse ponto de entrada é um comando —
então **`server-init` garante a camada CLI antes de renderizar a camada server**:

1. `hasCli := io.IsDir("sandbox/internal/cli")`. Se for `false`, `ServerInitInternal` chama
   `cli_init.CliInitInternal(deps, io, path)` sobre o **mesmo** SmartIO aberto (ações compõem
   compartilhando um `*SmartIO`, como manda o `CLAUDE.md`) — sem `Persist` intermediário e sem
   `build` intermediário.
2. Renderiza o grupo `server`.
3. Escreve `sandbox/internal/commands/start_server/{entries.yaml,handler.go}` a partir de
   `assets/templates/start_server_*` com `utils.RenderTemplateToDest`, recusando sobrescrever um
   `start_server/` já existente (mesma regra de `add-route`).
4. Persiste e roda o `build`, que gera o `entries.go` do `start-server` e o `climain.go` já com
   ele no despacho.

`server-purge` é o inverso e apaga `sandbox/internal/commands/start_server/` junto; **não** roda
`cli-purge` — a camada CLI, uma vez instalada, é do projeto.

### Comando `start-server`

Único ponto de entrada do servidor pela CLI. Handler à mão (scaffold), curto por construção: lê
as flags e chama `server.ServerMain` — o mesmo que `binds/server.go` expõe como
`sandbox.Server.Serve`.

```yaml
# sandbox/internal/commands/start_server/entries.yaml
identifiers: ["start-server"]
category: Server
help: Starts the http server
long-description: |
  Opens the port and serves every route declared under
  sandbox/internal/routes, until the process is stopped.
examples:
  - "start-server"
  - "start-server --addr :3000"
flags:
  - name: addr
    identifiers: ["--addr"]
    description: the address the server listens on
    examples:
      - "start-server --addr :3000"
    type: string
    default: ":8080"
  - name: read_timeout_ms
    identifiers: ["--read-timeout-ms"]
    description: how long a request has to arrive, in milliseconds
    examples:
      - "start-server --read-timeout-ms 30000"
    type: int
    default: "10000"
  - name: write_timeout_ms
    identifiers: ["--write-timeout-ms"]
    description: how long a response has to be written, in milliseconds
    examples:
      - "start-server --write-timeout-ms 30000"
    type: int
    default: "10000"
```

```go
// sandbox/internal/commands/start_server/handler.go
func CommandHandler(deps *deps.Deps, entries *Entries) int {
	err := server.ServerMain(deps, api.ServeProps{
		Addr:           entries.Addr,
		ReadTimeoutMs:  entries.ReadTimeoutMs,
		WriteTimeoutMs: entries.WriteTimeoutMs,
	})
	if err != nil {
		deps.Std.Error("server stopped: %s \n", err.Error())
		return api.ExitFailure
	}
	return api.ExitOk
}
```

`Serve` bloqueia; o handler só retorna quando o servidor cai. Como todo `handler.go`, é escrito
uma vez e nunca reescrito por `build`.

Em `sandbox/api/actions.go` (doc comment em cada campo) e `sandbox/binds/actions.go`:

```go
type RouteProps struct {
	Path, Route, Method, Help, Category, LongDescription string
	Hidden, Visible bool
	Examples []string
}

type RouteFieldProps struct {
	Path, Route, Name, Identifier, In, Description, Type, Default, Min, Max string
	Examples []string
	Required, Array bool
	Position int
	Format, Pattern string   // só para --in body
}

// em Actions:
ServerInit  func(path string) error
ServerPurge func(path string) error
AddRoute    func(path string, name string, method string, trigger string, help string, category string) error
RemoveRoute func(path string, name string) error
SetRoute    func(props RouteProps) error
AddSegment      func(props RouteFieldProps) error
RemoveSegment   func(path string, route string, name string) error
AddHeader       func(props RouteFieldProps) error
RemoveHeader    func(path string, route string, name string) error
AddParam        func(props RouteFieldProps) error
RemoveParam     func(path string, route string, name string) error
SetBody         func(props RouteBodyProps) error
AddBodyField    func(props RouteBodyFieldProps) error
RemoveBodyField func(path string, route string, name string) error
```

`trigger` (em `AddRoute`) e `Identifier` (em `RouteFieldProps`) são normalizados para começar
com `/`: `users` e `/users` produzem o mesmo `route.yaml`, e `/` interno ou final é erro.

---

## 8. Fase 7 — `verify`

`sandbox/internal/actions/verify/check_routes.go` com `CheckRoutes(deps, io) []string`, chamado
em `verify_internal.go`. Nenhuma escrita, uma string por violação:

- todo `routes/<name>/` tem `route.yaml`, `entries.go` e `handler.go`
- `route.yaml` parseia por `routeconf`; `method` é verbo conhecido
- `paths` não vazio, ao menos um `identifier`, cada item com `identifier` **ou** `name` (nunca
  os dois); todo `identifier` começa com `/`, nenhum é vazio, e só o da raiz (`/`) não tem nada
  depois da barra — `/` interno ou final é violação
- segmento capturado: `required` só `true`, sem `array` nem `default`
- nenhum `name` repetido na mesma origem; `name` em mais de uma origem tem o mesmo `type`
- nenhum par (`method`, padrão derivado de `paths`) repetido entre rotas
- `json-schema` só com `body.type: json`, e só com chaves do subconjunto
- `required`/`default` mutuamente exclusivos; `required` proibido em `boolean`
- `handler.go` exporta `RouteHandler` com a assinatura canônica
- projeto sem `sandbox/internal/routes/` -> nenhuma violação (padrão de `CheckStructure`)

---

## 9. Fase 8 — docs, structure, themes

| Arquivo | Edição |
|---|---|
| `assets/server/docs/RouteYaml/` | toda chave do `route.yaml` (irmão de `assets/all/docs/EntriesYaml/`) |
| `assets/server/docs/Routes/` | tabela de `{{range .RouteDocs}}` (irmão de `docs/Commands`) |
| `assets/server/docs/ServerUsage/` | subir o servidor; ciclo `add-route` -> editores -> `build` |
| `assets/all/docs/Rules/doc.md` | seção `## Routes` sob `{{ if .HasServer }}` |
| `assets/all/docs/GeneratedFiles/doc.md` | bloco `{{- if .HasServer }}`: `api/server.go`, `binds/server.go`, `internal/server/*.go`, `routes/<name>/entries.go` (always), `route.yaml`/`handler.go` (once) |
| `assets/all/docs/Workflow/doc.md` | `## Add the server layer`, `## Change the route surface` |
| `assets/all/docs/DepList/doc.md` | linha do `serverdeps` |
| `AgnosConfig/themes.yaml` + `assets/start/AgnosConfig/themes.yaml` | tema `server-usage` / `ServerUsage` |
| `AgnosConfig/structure.yaml` | `internal/server/servermain.go` (gen), `internal/routes/<name>` (dir), `assets/server` (dir) |
| `docs/Contributing/doc.md` | espelhar o padrão novo **no mesmo commit** |

`README.md` e os `Index.md` saem sozinhos do `build`. Gerados no projeto-alvo (não versionados
como fonte): `sandbox/api/server.go`, `sandbox/binds/server.go`, `sandbox/internal/server/*.go`,
`routes/<name>/entries.go`, `docs/{RouteYaml,Routes,ServerUsage}/`.

---

## 10. Fase 9 — exemplos

Criados só por comando; `result.yaml` só por `exec-test`/`update-test`.

```bash
./release/bootstrap.bin add-cli-example server-init
./release/bootstrap.bin add-cli-example add-route
./release/bootstrap.bin add-lib-example server-route
```

- `cli/server-init` — `start` + `server-init`; copia `sandbox/api/server.go`,
  `sandbox/internal/server/`, `sandbox/internal/routes/health/` e
  `sandbox/internal/commands/start_server/` (prova de que o `cli-init` implícito rodou)
- `cli/add-route` — `add-route` + `add-segment` + `add-header` + `add-param` + `set-body` +
  `add-body-field`; copia o `route.yaml` e o `entries.go` da rota
- `lib/server-route` — mesmo resultado pela API, copiando o mesmo conjunto

---

## 11. Ordem de execução

Cada passo termina com o ciclo de bootstrap do `CLAUDE.md` (`build`, `verify`, `build -q` +
`git diff --quiet`); nunca rodar um `agnos` instalado neste repo.

1. `serverdeps` (contrato + adapter + deplist + espelho + `DepList`) — `dep-install serverdeps`
   já tem de compilar
2. `routeconf` + `docs/RouteYaml` — só parse/render
3. Grupo `assets/server/` com `api`, `binds`, `servermain.go` sem rotas e a rota `health`
4. Collectors + `generate_route_entries.go` + vars do `build_internal.go`
5. `routeio` (`ValidateSchema` + `WriteError`) e o `ReadBody` gerado em cada `entries.go`
6. Ações e comandos + `api/actions.go` + `binds/actions.go`
7. `check_routes.go` no `verify`
8. Docs, `themes.yaml`, `structure.yaml`, `Contributing`
9. Exemplos + `exec-test`
10. `version` em `AgnosConfig/project.yaml` e `agnos publish`

---

## 12. Pontos em aberto

- **`pattern` de JSON Schema** precisa de regex, que o sandbox não tem. Recomendação:
  acrescentar `MatchPattern` a `stringsdeps` (adapter sobre `regexp`) — é uma linha de contrato
  e resolve `format` também. Alternativa: `pattern` fora do subconjunto na v1.
- **`start-server` é sempre escrito**, porque `server-init` instala a camada CLI quando ela
  falta. Um projeto que quiser subir o servidor sem passar pela CLI continua podendo chamar
  `sandbox.Server.Serve(...)` direto — o comando é conveniência, não o único caminho.
- **Sem TLS na v1** — certificados entram depois como campos novos de `ServerProps`, sem quebrar
  o contrato.
