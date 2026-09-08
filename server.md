# Plano de implementação — mecânica de servidores

Camada `server`, espelho exato da camada `cli` que já existe. Nada aqui inventa um padrão
novo: cada arquivo abaixo é a instância *server* de um arquivo *cli* que já existe no repo
(regra "Every file is an instance of a pattern", `docs/Rules/doc.md`).

## 1. Correspondência com a camada CLI

| Conceito | CLI (existe) | Server (a criar) |
|---|---|---|
| Contrato de entrada externa | `sandbox/deps/argvdeps/` | `sandbox/deps/serverdeps/` |
| Adapter | `adapters/libs/verb/` | `adapters/libs/serverdeps/` |
| Superfície pública | `sandbox/api/cli.go` | `sandbox/api/server.go` |
| Bind | `sandbox/binds/cli.go` | `sandbox/binds/server.go` |
| Despacho gerado | `sandbox/internal/cli/climain.go` | `sandbox/internal/server/servermain.go` |
| Unidade declarada | `sandbox/internal/commands/<name>/` | `sandbox/internal/routes/<name>/` |
| Declaração (yaml, editada por comando) | `entries.yaml` | `route.yaml` |
| Struct tipada (gerada) | `entries.go` -> `Entries` | `entries.go` -> `Entries` |
| Handler (à mão) | `handler.go` -> `CommandHandler` | `handler.go` -> `RouteHandler` |
| Parsable | `parsables/commandconf/` | `parsables/routeconf/` |
| Grupo de assets | `assets/cli/` | `assets/server/` |
| Instalação / remoção da camada | `cli-init` / `cli-purge` | `server-init` / `server-purge` |
| Editores da declaração | `add-command`, `add-flag`, `add-arg`, `set-command` | `add-route`, `add-field`, `remove-field`, `set-route`, `remove-route` |
| Gatilho no build | `hasCli := io.IsDir("sandbox/internal/cli")` | `hasServer := io.IsDir("sandbox/internal/server")` |

Invariante herdada: **o despacho valida tudo que não é o corpo antes do handler rodar**.
Assim como um `CommandHandler` nunca retorna `api.ExitUsage`, um `RouteHandler` nunca *decide*
um `400`/`404`/`405`/`415` — quem responde isso é `servermain.go`.

O corpo é a exceção deliberada: ele não é lido nem parseado antes do handler. Um handler só
quer o corpo depois de aprovar o resto (autorização, tenant, estado), então `Entries` não traz
um campo `Body` pronto — traz o método **`ReadBody`**, que lê, valida e converte sob demanda,
devolvendo uma struct (`type: json` com schema), `[]byte` (`raw`) ou `string` (`text`). Falha
de corpo continua sem ser decisão do handler: `ReadBody` já responde o erro e devolve o status
para o handler propagar.

---

## 2. Fase 1 — dep `serverdeps`

Dep nomeado pelo contrato (`serverdeps`), campo `Deps.Serverdeps`, lib adapter `serverdeps`,
sobre `net/http` (stdlib -> **sem** entrada em `assets/depsversion.yaml`).
O dep é **0 opinativo**: não conhece rota, método, padrão nem `{param}`. Ele abre a porta,
aplica os timeouts e entrega toda requisição à única função que recebeu em
`ServerProps.Handler`. Casamento de rota, `404` e `405` são do `servermain.go` gerado — quem se
adapta ao dep é o sandbox, nunca o contrário. Sem `ServeMux`, sem `PathValue`: trocar `net/http`
por outra implementação não muda uma linha do sandbox.

| Arquivo | Papel |
|---|---|
| `assets/deplist/serverdeps/sandbox/deps/serverdeps/serverdeps.go` | template do contrato |
| `assets/deplist/serverdeps/adapters/libs/serverdeps/serverdeps.go` | template do adapter |
| `sandbox/deps/serverdeps/serverdeps.go` | cópia instalada neste repo (byte-a-byte com o template renderizado — **(verify)** `check_deplist.go`) |
| `adapters/libs/serverdeps/serverdeps.go` | idem |

Forma do contrato (construção por chamada, igual a `requestdeps`; toda declaração exportada
leva doc comment, senão `verify` falha):

```go
package serverdeps

type Lib struct {
	NewServer func(props ServerProps) Server
}

type ServerProps struct {
	Addr           string // ":8080"
	ReadTimeoutMs  int
	WriteTimeoutMs int
	Handler        func(Request, Response) // chamada uma vez por requisição, qualquer que seja
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

Nada de `time.Time`, `io.Reader` ou tipos de `net/http` atravessa o contrato — só builtins
(mesma razão do commit `71cc554`, timeouts em `int` de milissegundos).

Editar também: `assets/all/docs/DepList/doc.md` (uma linha na tabela) e
`assets/depsversion.yaml` (nenhuma — registrar na descrição que é stdlib).

---

## 3. Fase 2 — `route.yaml`, a declaração

`sandbox/internal/routes/<name>/route.yaml`. Escrito por `add-route` e reescrito por
`add-field` / `remove-field` / `set-route` — **nunca à mão**, exatamente como `entries.yaml`.
Diretório em snake_case para rota kebab-case (`get-user` -> `get_user/`).

A declaração é lida de fora para dentro: **`paths`** diz qual URL casa, **`headers`** e
**`params`** dizem o que é lido dela, **`body`** diz o que vem no corpo. Não existe chave
`path` escrita à mão — o caminho da rota é a concatenação, na ordem, dos segmentos de `paths`,
e o collector é quem o deriva.

Um segmento de `paths` é de um de dois tipos, distinguidos pela chave presente:

| Segmento | Chave | Papel |
|---|---|---|
| **trigger** | `identifier` | literal na URL. O primeiro `identifier` da rota é o gatilho que identifica o handler — é ele que nomeia a rota em `docs/Routes` |
| **captura** | `name` | segmento variável, vira campo de `Entries` e chega ao handler já convertido |

```yaml
# sandbox/internal/routes/article/route.yaml -> GET /articles/{article-name}
method: GET
paths:
  - identifier: "articles"
  - name: "article-name"
    type: string
    required: true
category: Articles
help: Read one article by its slug
examples:
  - "curl -X GET http://[IP_ADDRESS]/articles/python-beginners-guide"
```

O mesmo `name` pode ser declarado em mais de uma origem — o campo de `Entries` é único e, em
runtime, vence a **primeira origem declarada** que trouxer valor (`paths`, depois `headers`,
depois `params`, na ordem do arquivo). `required` é satisfeito por qualquer uma delas; o
tipo tem de ser o mesmo em todas (**(verify)**).

```yaml
# sandbox/internal/routes/get_user/route.yaml -> GET /get-user-by-username
method: GET
paths:
  - identifier: "get-user-by-username"
headers:
  - name: "username"
    type: string
    required: true
params:
  # username está declarado nas duas origens: vale o header, que vem primeiro,
  # e a query só responde quando o header não veio
  - name: "username"
    type: string
    required: true
examples:
  - "curl -X GET http://[IP_ADDRESS]/get-user-by-username -H 'username: mateus'"
```

Exemplo completo, com corpo e schema:

```yaml
# sandbox/internal/routes/create_user/route.yaml -> POST /users/{tenant}/create
method: POST
paths:
  - identifier: "users"
  - name: "tenant"
    type: string
    required: true
  - identifier: "create"
category: Users
help: Create a user under a tenant
long-description: |
  Cria um usuário. Falha com 409 quando o e-mail já existe.
examples:
  - "curl -X POST localhost:8080/users/acme/create -d '{\"email\":\"a@b.c\",\"age\":30}'"
hidden: false
headers:
  - name: authorization
    description: bearer token
    type: string
    required: true
  - name: x-trace-id
    type: string
    default: ""
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
      email:
        type: string
        format: email
        minLength: 3
        maxLength: 254
      age:
        type: integer
        minimum: 0
        maximum: 130
      tags:
        type: array
        items: { type: string }
        minItems: 0
      address:
        type: object
        required: [city]
        properties:
          city: { type: string }
          zip:  { type: string, pattern: "^[0-9]{5}$" }
```

### Chaves de rota

| Chave | Efeito |
|---|---|
| `paths` | sequência de segmentos da URL, na ordem. Obrigatória e não vazia; o primeiro `identifier` é o gatilho do handler |
| `method` | `GET`/`POST`/`PUT`/`PATCH`/`DELETE`/`HEAD`/`OPTIONS`. Opcional, default `GET` |
| `category`, `help`, `long-description`, `examples`, `hidden` | idêntico a `entries.yaml`; alimenta `docs/Routes` |
| `headers` | sequência de campos lidos do cabeçalho |
| `params` | sequência de campos lidos da query string |
| `body` | descrição do corpo (um objeto, não sequência) |

### Chaves de segmento (`paths`)

| Chave | Efeito |
|---|---|
| `identifier` | literal do segmento. Exclui `name` e todas as chaves de campo no mesmo item |
| `name` | nome do segmento capturado; vira posição variável no casamento e campo de `Entries` |
| `type` | `string`\|`boolean`\|`int`\|`float`. Default `string` |
| `description`, `examples` | documentação, alimenta `docs/Routes` |
| `min`, `max` | mesmos limites de `entries.yaml` |
| `required` | sempre `true` num segmento capturado — `false` é erro de **(verify)** |

`array` e `default` são proibidos num segmento: a URL casa ou não casa.

### Chaves de campo (`headers`, `params`)

Mesmas de `entries.yaml` — `name`, `description`, `examples`, `type` (`string`|`boolean`|
`int`|`float`), `default`, `required`, `array`, `min`, `max`. O `name` **é** a grafia externa:
nome do header (casado sem diferenciar maiúsculas) ou chave na query string. `array: true` só
em `params`, e lê todas as ocorrências da chave.

### Chaves de `body`

| Chave | Efeito |
|---|---|
| `type` | `none` (padrão), `raw` (`[]byte`), `text` (`string`), `json` |
| `required` | corpo ausente/vazio vira `400` |
| `max-bytes` | corte de leitura; excedido vira `413`. Default 1 MiB |
| `content-type` | exigido no header; divergência vira `415`. Default `application/json` para `type: json` |
| `json-schema` | subconjunto de JSON Schema, só quando `type: json` |

Só `content-type` é checado pelo despacho (é um header, não custa ler o corpo). `required`,
`max-bytes` e `json-schema` são aplicados dentro de `ReadBody`, na primeira chamada do handler
— um `Content-Length` já maior que `max-bytes` é a única antecipação, e responde `413` antes de
chamar o handler.

### Subconjunto de JSON Schema suportado

`type` (`object`/`array`/`string`/`integer`/`number`/`boolean`/`null`), `properties`,
`required`, `additionalProperties` (bool), `items`, `enum`, `const`, `minimum`, `maximum`,
`exclusiveMinimum`, `exclusiveMaximum`, `minLength`, `maxLength`, `pattern`, `minItems`,
`maxItems`, `uniqueItems`, `format` (`email`, `uuid`, `date-time`, `uri`), `nullable`.
Fora do subconjunto (`$ref`, `oneOf`, `allOf`, `anyOf`, `patternProperties`) -> erro de
`verify`/`build`, não silêncio.

---

## 4. Fase 3 — parsable `routeconf`

`sandbox/internal/parsables/routeconf/` — cinco arquivos, mesmos nomes e mesma ordem que
`commandconf/`:

| Arquivo | Conteúdo |
|---|---|
| `sandbox/internal/parsables/routeconf/api.go` | `RouteConf` (com `Paths []Segment`), `Segment` (`Identifier` ou `Field`), `Field`, `Body`, `Schema` |
| `sandbox/internal/parsables/routeconf/new.go` | `New(deps, content) (*RouteConf, error)` |
| `sandbox/internal/parsables/routeconf/new_empty.go` | `NewEmpty(deps) *RouteConf` |
| `sandbox/internal/parsables/routeconf/bind_methods.go` | `BindMethods` (liga `Render`) |
| `sandbox/internal/parsables/routeconf/render.go` | `Render` -> volta ao `route.yaml` canônico |

`Segment` é `{Identifier string; Field *Field}` — exatamente um dos dois preenchido; o caminho
canônico (`RouteConf.Pattern()`) é `/` + os segmentos na ordem, literal ou `{name}`, usado em
`docs/Routes` e nas mensagens — o casamento em si compara segmento a segmento, não o texto.
`Schema` é uma árvore recursiva (`Type`, `Properties []SchemaProperty`, `Items *Schema`,
`Required []string`, bounds com o par valor/`Has…` como em `Field.Min`/`HasMin`).
`Render` reemite as chaves na ordem canônica — determinismo/idempotência.

---

## 5. Fase 4 — a camada server (grupo `assets/server/`)

Renderizado por todo `build` quando `sandbox/internal/server/` existe. Todo arquivo do grupo é
**sempre reescrito**.

| Template | Destino no projeto |
|---|---|
| `assets/server/sandbox/api/server.go` | `sandbox/api/server.go` |
| `assets/server/sandbox/binds/server.go` | `sandbox/binds/server.go` |
| `assets/server/sandbox/internal/server/servermain.go` | `sandbox/internal/server/servermain.go` |
| `assets/server/sandbox/internal/routeio/jsonschema.go` | `sandbox/internal/routeio/jsonschema.go` |
| `assets/server/sandbox/internal/routeio/write_error.go` | `sandbox/internal/routeio/write_error.go` |
| `assets/server/sandbox/internal/routes/health/route.yaml` | `sandbox/internal/routes/health/route.yaml` |
| `assets/server/sandbox/internal/routes/health/handler.go` | `sandbox/internal/routes/health/handler.go` |
| `assets/server/docs/RouteYaml/{doc.md,props.yaml}` | `docs/RouteYaml/` |
| `assets/server/docs/Routes/{doc.md,props.yaml}` | `docs/Routes/` (renderizado das rotas, como `docs/Commands`) |
| `assets/server/docs/ServerUsage/{doc.md,props.yaml}` | `docs/ServerUsage/` |

`health` é a rota embutida — o análogo do comando `version`, para que um servidor recém-criado
já responda algo.

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
	StatusOk                 = 200
	StatusCreated            = 201
	StatusNoContent          = 204
	StatusBadRequest         = 400
	StatusNotFound           = 404
	StatusMethodNotAllowed   = 405
	StatusConflict           = 409
	StatusPayloadTooLarge    = 413
	StatusUnsupportedMedia   = 415
	StatusFailure            = 500
)
```

`sandbox/api/sandbox.go` ganha o campo `Server` sozinho (já é gerado com um campo por arquivo
de `api/`), e `sandbox/new.go` chama `binds.ServerBind` sozinho (um `<X>Bind` por arquivo de
`binds/`). Nada a editar à mão nesses dois.

### `servermain.go` (gerado)

Uma única função entregue ao dep, e todo o despacho abaixo dela — o dep não sabe que rota
existe. `ServerMain` não tem `range`; o `{{range .Routes}}` está no `dispatch`:

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

mais um `match<GoName>(segments []string) bool` por rota — comprimento mais os `identifier`
comparados por posição, as posições de captura aceitando qualquer valor —, e um
`handle<GoName>` por rota, que: lê e converte os segmentos capturados de `paths` (por índice em
`segments`, não por nome) ->
headers -> params, na ordem de declaração e com a primeira origem que traz valor vencendo em
nomes repetidos; aplica defaults; checa `required`, `min`/`max`, `array`; valida `content-type`
e um `Content-Length` acima de `max-bytes`; monta `Entries` — guardando em `Entries.Request` a
requisição de onde `ReadBody` vai ler; e só então chama
`routes_<name>.RouteHandler(deps, &entries, res)`. **O corpo não é tocado aqui**: nada é lido do
socket enquanto o handler não chamar `ReadBody`. Qualquer falha responde antes, em JSON
(`{"error": "...", "field": "..."}`), com o status da tabela abaixo, e registra em
`deps.Std.Log`. Os helpers compartilhados (`readHeader`, `parseIntValue`, `bindQueryArray`)
ficam no rodapé do arquivo, junto de `splitPath` (`deps.Stringsdeps.Split` do caminho por
`"/"`, descartando os vazios das pontas), exatamente como os helpers de argv em `climain.go`; o
`writeError` é o de `routeio`, compartilhado com o `ReadBody` de cada rota.

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

Pacote próprio, e não um arquivo dentro de `sandbox/internal/server/`, porque quem chama o
validador agora é o `ReadBody` de cada rota: `servermain.go` importa todo
`sandbox/internal/routes/<name>/`, então uma rota não pode importar `internal/server` de volta.
`routeio` não importa nenhum dos dois e é importado pelos dois.

`jsonschema.go` — validador puro do subconjunto acima, sobre `deps.Serializables.ParseJson` —
sem stdlib, sem regexp externo (`pattern` limitado às âncoras/classes que o validador
implementa, ou delegado a um `Stringsdeps.MatchPattern` a acrescentar no contrato):

```go
func ValidateSchema(deps *deps.Deps, schema_json string, body []byte) (*serializables.SerializibleObject, string, bool)
```
(o objeto parseado, a mensagem da primeira violação, e se passou).

`write_error.go` — a única forma de escrever uma falha, usada pelo despacho e por todo
`ReadBody`, para que os dois lados respondam o mesmo JSON:

```go
func WriteError(deps *deps.Deps, response serverdeps.Response, status int, field string, message string) int
```
(devolve o próprio `status`, para o chamador retornar em uma linha).

---

## 6. Fase 5 — geração por rota, no `build`

| Arquivo | Papel |
|---|---|
| `sandbox/internal/actions/build/collect_routes.go` | `CollectRoutes(deps, io)`: lê todo `sandbox/internal/routes/<name>/route.yaml` via `routeconf`, devolve `[]map[string]any` (`Name`, `GoName`, `Method`, `Trigger`, `Path` — derivado de `paths` —, `Segments`, `Headers`, `Params`, `Body`, `SchemaJson`, `BodyStructs`, `HasBody`) — cópia de `collect_commands.go` |
| `sandbox/internal/actions/build/collect_route_docs.go` | alimenta `docs/Routes` — cópia de `collect_command_docs.go` |
| `sandbox/internal/actions/build/generate_route_entries.go` | renderiza `assets/templates/route_entries.go` uma vez por rota em `sandbox/internal/routes/<name>/entries.go` |
| `assets/templates/route_entries.go` | template do `Entries` de uma rota |
| `assets/templates/route_route.yaml` | scaffold de `route.yaml` (usado por `add-route`) |
| `assets/templates/route_handler.go` | scaffold de `handler.go` (usado por `add-route`) |
| `assets/templates/serve_entries.yaml` | scaffold do comando `serve` (escrito por `server-init` quando há CLI) |
| `assets/templates/serve_handler.go` | idem |

Edições em `sandbox/internal/actions/build/build_internal.go`:

- `hasServer := io.IsDir("sandbox/internal/server")`
- `routes, err := CollectRoutes(deps, io)` e `route_docs, err := CollectRouteDocs(...)`
- vars novas: `"HasServer"`, `"Routes"`, `"RouteDocs"`
- `GeneratedDocsGroups(hasCli, hasServer)` em `collect_generated_docs.go` passa a devolver
  `"server"` também
- ao final: `if hasServer { GenerateRouteEntries(...); utils.RenderGroup(deps, io, "server", vars) }`

### `entries.go` gerado de uma rota

```go
package create_user

type Entries struct {
	Tenant        string             // segmentos capturados de `paths`, na ordem declarada
	Authorization string             // headers, na ordem declarada
	XTraceId      string
	Page          int                // params, na ordem declarada
	Tag           []string
	Request       serverdeps.Request // preenchido pelo despacho; a fonte de `ReadBody`
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

// ReadBody lê, valida e converte o corpo na primeira chamada; o resultado fica em cache para
// as seguintes. Devolve api.StatusOk quando passou; em qualquer outro status a resposta de
// erro já foi escrita e o handler só propaga o retorno.
func (entries *Entries) ReadBody(deps *deps.Deps, response serverdeps.Response) (Body, int)
```

O corpo nunca é campo de `Entries` — é o retorno de `ReadBody`, que é gerado com o tipo de
`body.type`:

| `body.type` | Assinatura gerada |
|---|---|
| `none` (padrão) | nenhum `ReadBody` é gerado |
| `raw` | `ReadBody(deps, response) ([]byte, int)` |
| `text` | `ReadBody(deps, response) (string, int)` |
| `json` com `json-schema` de objeto | `ReadBody(deps, response) (Body, int)` |
| `json` sem `json-schema` | `ReadBody(deps, response) (*serializables.SerializibleObject, int)` |

Toda variante faz, na ordem: `Request.ReadBody(max-bytes)` (`413` se estourar), `required`
(`400` se vazio) e — só no caso `json` — `routeio.ValidateSchema` contra `EntriesSchema`
(`400` na primeira violação, com o campo reprovado em `"field"`). Uma segunda chamada devolve
o valor já lido, sem tocar no socket de novo.

Regras de nome: objeto aninhado -> `Body<Caminho>` (`Body` + `Address` -> `BodyAddress`);
item de array de objetos -> `Body<Caminho>Item`.

### `handler.go` (à mão, único arquivo escrito à mão de uma rota)

```go
func RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int
```

Retorna o status com que respondeu — o análogo do `int` de `CommandHandler`. Nunca *decide* um
`400`/`404`/`405`/`413`/`415`; o único caminho pelo qual devolve um deles é propagando o que
`ReadBody` já respondeu:

```go
func RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int {
	// checagens do handler primeiro — o corpo ainda não foi lido
	if !isAuthorized(deps, entries.Authorization) {
		return writeJson(response, api.StatusFailure, ...)
	}

	body, status := entries.ReadBody(deps, response)
	if status != api.StatusOk {
		return status
	}

	return writeJson(response, api.StatusCreated, createUser(deps, entries.Tenant, body))
}
```

Uma rota com `body.type: none` não tem `ReadBody` e nem essas três linhas.

---

## 7. Fase 6 — ações e comandos

Duas camadas por feature, como sempre: `<name>.go` (abre SmartIO, persiste, dispara `build`)
+ `<name>_internal.go` (lógica sobre SmartIO aberto).

| Ação (`sandbox/internal/actions/<dir>/`) | Comando (`sandbox/internal/commands/<dir>/`) | Faz |
|---|---|---|
| `server_init/` | `server_init/` -> `server-init` | instala `serverdeps`, `std`, `serializables`, `stringsdeps`; renderiza o grupo `server`; escreve o comando `serve` quando há CLI; roda `build` |
| `server_purge/` | `server_purge/` -> `server-purge` | remove os arquivos do grupo `server` + os diretórios `sandbox/internal/server` e `sandbox/internal/routes` inteiros (cópia de `cli_purge`) |
| `add_route/` | `add_route/` -> `add-route` | escreve `route.yaml` + `handler.go` de uma rota nova (recusa sobrescrever) |
| `remove_route/` | `remove_route/` -> `remove-route` | apaga o diretório da rota |
| `set_route/` | `set_route/` -> `set-route` | reescreve as chaves de nível de rota (`method`, `help`, `category`, `long-description`, `hidden`, `examples`) |
| `add_field/` | `add_field/` -> `add-field` | acrescenta um segmento ou campo, conforme `--in path\|header\|query\|body` |
| `remove_field/` | `remove_field/` -> `remove-field` | remove um campo declarado |

Um único par `add-field`/`remove-field` (e não um comando por origem) porque path, header,
query e body diferem só pela origem — `--in` é o seletor. `--in path` acrescenta um segmento
ao fim de `paths` (ou em `--position`): `--identifier <literal>` para um trigger,
`--name <nome>` para uma captura. Para `--in body`, `--name` aceita caminho pontuado
(`address.city`) e cria os objetos intermediários no `json-schema`.

Cada `sandbox/internal/commands/<dir>/` traz `entries.yaml` (à mão, via `add-command` +
`add-flag`/`add-arg`), `entries.go` (gerado) e `handler.go` (à mão).

Contratos a acrescentar em `sandbox/api/actions.go` (com doc comment em cada campo, senão
`verify` falha) e o bind correspondente em `sandbox/binds/actions.go`:

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
ServerInit   func(path string) error
ServerPurge  func(path string) error
AddRoute     func(path string, name string, method string, trigger string, help string, category string) error
RemoveRoute  func(path string, name string) error
SetRoute     func(props RouteProps) error
AddField     func(props RouteFieldProps) error
RemoveField  func(path string, route string, in string, name string) error
```

---

## 8. Fase 7 — `verify`

Novo arquivo `sandbox/internal/actions/verify/check_routes.go`, com `CheckRoutes(deps, io)
[]string`, chamado em `sandbox/internal/actions/verify/verify_internal.go` junto dos demais.
Regras (nenhuma escrita, uma string por violação):

- todo `sandbox/internal/routes/<name>/` tem `route.yaml`, `entries.go` e `handler.go`
- `route.yaml` parseia por `routeconf`; `method` é um verbo conhecido
- `paths` não é vazio, tem ao menos um `identifier`, e cada item traz `identifier` **ou**
  `name`, nunca os dois; nenhum `identifier` com `/` ou vazio
- num segmento capturado: `required` só `true`, nada de `array` nem `default`
- nenhum `name` repetido dentro da mesma origem, e um `name` declarado em mais de uma origem
  tem o mesmo `type` em todas
- nenhum par (`method`, padrão derivado de `paths`) repetido entre rotas
- `json-schema` só com `body.type: json`; só chaves do subconjunto suportado
- `required`/`default` mutuamente exclusivos; `required` proibido em `boolean` (mesmas regras
  de `entries.yaml`)
- `handler.go` exporta `RouteHandler` com a assinatura canônica
- projeto sem `sandbox/internal/routes/` -> nenhuma violação (mesmo padrão de `CheckStructure`)

---

## 9. Fase 8 — docs, structure, themes

| Arquivo | Edição |
|---|---|
| `assets/server/docs/RouteYaml/doc.md` + `props.yaml` | toda chave do `route.yaml` (irmão de `assets/all/docs/EntriesYaml/`) |
| `assets/server/docs/Routes/doc.md` + `props.yaml` | tabela renderizada de `{{range .RouteDocs}}` (irmão de `docs/Commands`) |
| `assets/server/docs/ServerUsage/doc.md` + `props.yaml` | subir o servidor, ciclo `add-route` -> `add-field` -> `build` |
| `assets/all/docs/Rules/doc.md` | seção `## Routes` sob `{{ if .HasServer }}` — o handler nunca decide status de validação (só propaga o de `ReadBody`), o corpo só é lido por `ReadBody`, `route.yaml` nunca é editado à mão, nomes canônicos |
| `assets/all/docs/GeneratedFiles/doc.md` | bloco `{{- if .HasServer }}` com as linhas de `sandbox/api/server.go`, `binds/server.go`, `internal/server/*.go`, `routes/<name>/entries.go` (always) e `route.yaml`/`handler.go` (once) |
| `assets/all/docs/Workflow/doc.md` | seções `## Add the server layer` e `## Change the route surface` |
| `assets/all/docs/DepList/doc.md` | linha do `serverdeps` |
| `AgnosConfig/themes.yaml` e `assets/start/AgnosConfig/themes.yaml` | tema novo `server-usage` / `ServerUsage` |
| `AgnosConfig/structure.yaml` | `sandbox/internal/server/servermain.go` (gen), `internal/routes/<name>` (dir), `assets/server` (dir) — `verify` reprova entrada cujo caminho não existe |
| `docs/Contributing/doc.md` | espelhar o padrão novo **no mesmo commit** |

`README.md` e os `Index.md` saem sozinhos do `build`.

---

## 10. Fase 9 — exemplos

Criados só por comando, nunca à mão; `result.yaml` só por `exec-test`/`update-test`.

```bash
./release/bootstrap.bin add-cli-example server-init
./release/bootstrap.bin add-cli-example add-route
./release/bootstrap.bin add-lib-example server-route
```

- `examples/cli/server-init/example.sh` — `start` + `server-init`, copia
  `sandbox/api/server.go`, `sandbox/internal/server/`, `sandbox/internal/routes/health/`
- `examples/cli/add-route/example.sh` — `add-route` + três `add-field` (`--in path`,
  `--in header`, `--in body`), copia o `route.yaml` e o `entries.go` da rota
- `examples/lib/server-route/example.go` — mesmo resultado pela API, copiando o mesmo conjunto
  (os dois lados de um mesmo nome têm de deixar a mesma árvore)

---

## 11. Onde cada arquivo mora — lista consolidada

```
assets/deplist/serverdeps/sandbox/deps/serverdeps/serverdeps.go      novo
assets/deplist/serverdeps/adapters/libs/serverdeps/serverdeps.go     novo
sandbox/deps/serverdeps/serverdeps.go                                novo (espelho)
adapters/libs/serverdeps/serverdeps.go                               novo (espelho)

assets/server/sandbox/api/server.go                                  novo
assets/server/sandbox/binds/server.go                                novo
assets/server/sandbox/internal/server/servermain.go                  novo
assets/server/sandbox/internal/routeio/jsonschema.go                 novo
assets/server/sandbox/internal/routeio/write_error.go                novo
assets/server/sandbox/internal/routes/health/route.yaml              novo
assets/server/sandbox/internal/routes/health/handler.go              novo
assets/server/docs/RouteYaml/{doc.md,props.yaml}                     novo
assets/server/docs/Routes/{doc.md,props.yaml}                        novo
assets/server/docs/ServerUsage/{doc.md,props.yaml}                   novo

assets/templates/route_entries.go                                    novo
assets/templates/route_route.yaml                                    novo
assets/templates/route_handler.go                                    novo
assets/templates/serve_entries.yaml                                  novo
assets/templates/serve_handler.go                                    novo

sandbox/internal/parsables/routeconf/{api,new,new_empty,bind_methods,render}.go   novo

sandbox/internal/actions/build/collect_routes.go                     novo
sandbox/internal/actions/build/collect_route_docs.go                 novo
sandbox/internal/actions/build/generate_route_entries.go             novo
sandbox/internal/actions/build/build_internal.go                     editado
sandbox/internal/actions/build/collect_generated_docs.go             editado

sandbox/internal/actions/{server_init,server_purge,add_route,remove_route,set_route,add_field,remove_field}/
    <name>.go + <name>_internal.go                                   novo
sandbox/internal/commands/{server_init,server_purge,add_route,remove_route,set_route,add_field,remove_field}/
    entries.yaml + entries.go + handler.go                           novo

sandbox/internal/actions/verify/check_routes.go                      novo
sandbox/internal/actions/verify/verify_internal.go                   editado
sandbox/api/actions.go                                               editado
sandbox/binds/actions.go                                             editado

assets/all/docs/{Rules,GeneratedFiles,Workflow,DepList}/doc.md       editado
AgnosConfig/{themes.yaml,structure.yaml}                             editado
assets/start/AgnosConfig/themes.yaml                                 editado
docs/Contributing/doc.md                                             editado

examples/cli/server-init/, examples/cli/add-route/,
examples/lib/server-route/                                           via add-*-example
```

Gerados no projeto-alvo (não versionados como fonte): `sandbox/api/server.go`,
`sandbox/binds/server.go`, `sandbox/internal/server/*.go`,
`sandbox/internal/routes/<name>/entries.go`, `docs/{RouteYaml,Routes,ServerUsage}/`.

---

## 12. Ordem de execução

Cada passo termina com o ciclo de bootstrap; nunca rodar um `agnos` instalado neste repo.

```bash
go build -o release/bootstrap.bin ./cmd/main
./release/bootstrap.bin build
./release/bootstrap.bin verify
./release/bootstrap.bin build -q && git diff --quiet && echo idempotent
```

1. `serverdeps` (contrato + adapter + deplist + espelho + linha em `DepList`) — `dep-install
   serverdeps` já tem de compilar
2. `routeconf` + `docs/RouteYaml` — só parse/render, sem geração ainda
3. Grupo `assets/server/` com `api`, `binds`, `servermain.go` sem rotas e a rota `health`
4. Collectors + `generate_route_entries.go` + vars do `build_internal.go`
5. `routeio` (`ValidateSchema` + `WriteError`) e o `ReadBody` gerado em cada `entries.go`
6. Ações e comandos (`server-init`, `add-route`, `add-field`, …) + `api/actions.go` +
   `binds/actions.go`
7. `check_routes.go` no `verify`
8. Docs, `themes.yaml`, `structure.yaml`, `Contributing`
9. Exemplos + `exec-test`
10. `version` em `AgnosConfig/project.yaml` e `agnos publish`

---

## 13. Decisões e pontos em aberto

- **Dep 0 opinativo; o casamento de rota é do sandbox.** Um dep não decide nada do domínio de
  quem o usa — `serverdeps` só sabe abrir porta, receber requisição e escrever resposta, e
  recebe uma única `Handler` em `ServerProps`. Rota, método, `{param}`, `404` e `405` são
  política da aplicação e vivem no `servermain.go` gerado: é o sandbox que se adapta ao dep. O
  custo é um `match<GoName>` gerado por rota; o ganho é que o contrato não vaza a semântica de
  `ServeMux` (padrões, precedência, `PathValue`) e um adapter sobre outra implementação
  continua sendo uma troca de arquivo.
- **Schema validado em runtime, struct gerada em build.** A struct sai do `json-schema` no
  `build` (tipagem estática no handler) e o mesmo schema vai canonizado para
  `EntriesSchema`, validado por `routeio.ValidateSchema` a cada request. Uma só declaração,
  dois usos.
- **Corpo lido sob demanda, tudo o mais no despacho.** Path, headers e params são baratos e
  identificam a requisição, então o despacho os resolve antes do handler. O corpo é caro e só
  interessa depois que o handler aprovou o resto — uma requisição rejeitada por autorização não
  paga a leitura nem o parse. Daí `ReadBody` em vez de um campo `Body`, e daí `routeio` como
  pacote separado (uma rota não pode importar `internal/server`, que a importa).
- **`pattern` de JSON Schema** precisa de regex, que o sandbox não tem. Ou se acrescenta
  `MatchPattern` a `stringsdeps` (adapter sobre `regexp`), ou `pattern` fica fora do
  subconjunto na primeira versão. Recomendação: acrescentar ao `stringsdeps`, é uma linha de
  contrato e resolve `format` também.
- **`serve` como comando** só é escrito quando o projeto tem camada CLI; um projeto lib sobe o
  servidor por `sandbox.Server.Serve(...)`.
- **Sem TLS na primeira versão** — `Addr` e timeouts bastam; certificados entram depois como
  campos novos de `ServerProps`, sem quebrar o contrato.
