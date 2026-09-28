# Nova CLI — alinhar o sistema de comandos ao sistema de rotas

Objetivo: o `sandbox-cli` passa a espelhar o `sandbox-server` arquivo por arquivo, verbo por verbo,
para que as extensões sigam um só padrão. Rotas são a referência; onde a cli difere sem motivo, ela
muda. Onde a diferença é do domínio (argv ≠ http), o documento diz por quê.

## 1. Diagnóstico — o que difere hoje

### 1.1 Arquivos de uma unidade

| Papel | Rota (`sandbox/internal/routeslist/<name>/`) | Comando (`sandbox/internal/commands/<snake>/`) |
|---|---|---|
| declaração | `route.yaml` | `entries.yaml` |
| declaração em Go | `new.go` sobre a base genérica `generated/server/route.NewRoute` | `new.go` sobre `api.NewCommand()` — sem base genérica |
| valores tipados | `entries.go` → `type Entries struct` com tag `id`, gerado | **não existe**; lê por string: `command.GetString("path")` |
| lógica escrita à mão | `InternalPureHandler.go` | `handler.go` |
| assinatura | `InternalPureHandler(sandbox, props *api.RouteProps, entries *Entries, response *serverdeps.Response) error` | `CommandHandler(sandbox, command *api.Command) int` |
| estado compartilhado | `props *api.RouteProps`, tipado pelo projeto em `sandbox/api/routeprops.go` (escrito uma vez) | nenhum |
| falha | `return routeio.Fail(...)` → `routeio.Raise` → um dos 8 `sandbox/internal/server/errors/handle_*.go` (escritos uma vez, do projeto) | cada handler faz `Std.Error(...)` + `return api.ExitFailure`; mensagens de uso fixas dentro de `climain.go` |
| pacote io gerado | `generated/routeio/` (Fail, Raise, values, respond…) | nenhum; tudo inline em `climain.go` (430 linhas) |
| binder | `generated/server/route/RequestHandler.go` (genérico, via `Reflectdeps`) + `IsActionable.go` | `bindFlag` / `bindArg` dentro de `climain.go` |
| doc por unidade | `docs/Routes/<route>.md` | `docs/Commands/<cmd>.md` ✅ já alinhado |
| unidade embutida | `health` | `help`, `version` |
| `verify` | `check_routes.go` (assinatura, nome do arquivo, chaves) | **nenhum `check_commands.go`** |

**Colisão de nomes:** `entries.yaml` (declaração da cli) e `entries.go` (struct tipado da rota) são
coisas diferentes com o mesmo radical. Um LLM que leu uma extensão erra a outra.

### 1.2 Verbos de edição

| Operação | Rota | Comando | Falta na cli |
|---|---|---|---|
| criar unidade | `add-route` (tudo opcional, defaults) | `add-command` (`--help` e `--category` **obrigatórios**) | defaults |
| editar unidade | `set-route` com `--clear`, `--method` substitui a lista | `set-command` sem `--clear`, `--identifier` só acrescenta | `--clear`, substituição |
| remover unidade | `remove-route` | `remove-command` | — |
| renomear unidade | `rename-route` (move dir + reescreve `package`) | — | `rename-command` |
| listar | `list-routes` | — (só `help`) | `list-commands` |
| mostrar declaração | `show-route` (árvore, não escreve, não roda build) | — | `show-command` |
| simular entrada | `explain-route <method> <path>` | — | `explain-command <argv…>` |
| adicionar entrada | `add-path`, `add-parameter`, `add-body-field` | `add-flag`, `add-arg` | — |
| editar entrada | `set-path`, `set-parameter`, `set-body-field` (`--rename`, `--clear`, mesmo construtor do `add-`) | — | `set-flag`, `set-arg` |
| remover entrada | `remove-path`, `remove-parameter`, `remove-body-field` | `remove-flag`, `remove-arg` | — |
| importar em lote | `import-body` | — | fora de escopo (ver §6) |
| ordenação | `rebalance-routes` | — | n/a (argv não é cadeia) |

### 1.3 Vocabulário de tipos

| | Rota | Comando |
|---|---|---|
| inteiro | `integer` | `int` |
| real | `number` | `float` |
| lista | `string-array`, `integer-array` (tipo) | `--array` (flag) |
| validação | `--min/--max`, `--enum`, `--pattern`, `--format`, `trigger` | só `--min/--max` |

## 2. Alvo — a unidade de comando nova

```
sandbox/internal/commands/<snake>/
  command.yaml              # declaração (era entries.yaml) — espelho de route.yaml
  new.go                    # gerado: NewCommand sobre generated/cli/command.NewCommand
  entries.go                # gerado: type Entries struct, um campo por flag/arg, tag `id`
  InternalPureHandler.go    # à mão (era handler.go)
```

```go
// InternalPureHandler.go — escrito uma vez pelo add-command, depois do projeto
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries) error {
	sandbox.Deps.Std.Printf("my-feature called\n")
	return nil
}
```

```go
// entries.go — gerado
type Entries struct {
	Path  string   `id:"path"`
	Quiet bool     `id:"quiet"`
	Name  string   `id:"name"`
	Tags  []string `id:"tags"`
}
```

Regras, copiadas das rotas:

- O handler retorna `error`, nunca um código. `nil` → `ExitOk`. Recusar é `return cliio.Fail(sandbox, api.ExitFailure, "<id>", "<mensagem>")`. Qualquer outro `error` → `handle_failure.go`. Isso elimina o bloco `if err != nil { Std.Error; return ExitFailure }` repetido em ~80 handlers deste repo.
- `ExitUsage` continua exclusivo da dispatch (regra atual do CLAUDE.md).
- Campo de `Entries` = nome Go exportado do `id` (`out-file` → `OutFile`), regra de `RouteEntryId`. Nunca `Entries` com campo `Command` ou `Props` (reservados, como `FullRoute`/`Body`).
- O binder genérico preenche `Entries` pela tag `id` via `Deps.Reflectdeps`, como `RequestHandler.go`.
- `api.Command.Items` / `GetString` ficam **só** durante a migração (§5) e depois saem.

### 2.1 `command.yaml` — mesmo formato de `route.yaml`

```yaml
identifiers: [my-feature]
category: Commands
help: Does the thing
long-description: ""
hidden: false
examples:
  - my-feature ./dir --tag a
flags:
  - id: path
    identifiers: [--path]
    type: string
    default: .
    description: the dir holding the project
args:
  - id: name
    type: string
    required: true
```

Mudanças de esquema:

| Hoje | Novo | Motivo |
|---|---|---|
| `name:` numa flag/arg | `id:` | igual a `paths[].id` / `parameters[].id` |
| `type: int` / `float` | `integer` / `number` (aceitar `int`/`float` como alias na entrada, gravar só o canônico) | vocabulário único |
| — | `enum:`, `pattern:` | mesmas palavras de `add-body-field`; validados na dispatch, `ExitUsage` |
| arquivo `entries.yaml` | `command.yaml` | acaba a colisão com `entries.go`; `docs/EntriesYaml/` vira `docs/CommandYaml/`, par de `RouteYaml` |

## 3. Camadas geradas — espelho do server

| Server | Cli nova | Conteúdo |
|---|---|---|
| `generated/server/server/servermain.go` | `generated/cli/cli/climain.go` | só dispatch: acha o comando pelo identificador, `--help`, chama o handler genérico |
| `generated/server/server/new.go` | `generated/cli/cli/new.go` | lista `Cli.Commands`, preenche `Cli.Fail` |
| `generated/server/route/new.go` | `generated/cli/command/new.go` | base genérica `NewCommand(sandbox)` |
| `generated/server/route/RequestHandler.go` | `generated/cli/command/CommandHandler.go` | bind argv → `Entries` via Reflectdeps, chama `InternalPureHandler` |
| `generated/server/route/IsActionable.go` | `generated/cli/command/IsActionable.go` | `answersTo`; lido também por `utils/command_match.go` para `explain-command` |
| `generated/routeio/` | `generated/cliio/` | `Fail`, `FailWithCause`, `Raise`, `values.go` (parseValue, inRange, defaultValue saem do climain) |
| `sandbox/internal/server/errors/handle_*.go` (8, escritos uma vez) | `sandbox/internal/cli/errors/handle_*.go` (escritos uma vez) | `handle_unknown_command`, `handle_unknown_flag`, `handle_missing_value`, `handle_bad_value`, `handle_out_of_range`, `handle_unexpected_arg`, `handle_failure` |
| `sandbox/api/routeprops.go` (escrito uma vez) | `sandbox/api/commandprops.go` (escrito uma vez) | estado do processo que o projeto tipa |
| `api.RouteFailure` | `api.CommandFailure{ExitCode, Field, Message, Cause}` | |
| `api.Server.Fail` | `api.Cli.Fail` | mesma razão: comando não importa `generated/cli/cli` |

Assinatura dos `handle_*.go` da cli, par de `(sandbox, route, response) error`:
`(sandbox *api.Sandbox, command *api.Command) int` — o único ponto que devolve código de saída.

Como no server: nenhuma mensagem de erro fica escrita dentro da dispatch; uma falha sem mensagem
usa o texto do `handle_*.go` do projeto, e editar esse arquivo muda o que a cli diz.

## 4. Verbos — o que criar e mudar

### 4.1 Novos

| Comando | Espelho de | Faz |
|---|---|---|
| `set-flag <name> --command <c>` | `set-path` / `set-parameter` | `--rename`, todos os campos de `add-flag`, `--identifier` **substitui** a lista, `--clear <key>` (repetível: `default`, `required`, `min`, `max`, `description`, `examples`, `enum`, `pattern`), reconstrói pelo mesmo `utils.NewField` |
| `set-arg <name> --command <c>` | idem | idem, mais `--position` para mover |
| `rename-command <command> <name>` | `rename-route` | move o dir, reescreve a cláusula `package` dos `.go` à mão, troca o identificador principal, roda build. Recusa `help`/`version` |
| `list-commands` | `list-routes` | uma linha por comando: categoria, identificadores, nº de flags/args, nome. Não escreve |
| `show-command <command>` | `show-route` | árvore do `command.yaml`: linha de uso, flags, args com cada chave. Não escreve, não roda build |
| `explain-command <argv…>` | `explain-route` | percorre a dispatch sem rodar handler: qual comando casa, cada flag/arg ligado ou por que falha, e o `handle_*` que responderia. Não escreve |

Todo novo verbo: pasta de action com `<name>.go` + `<name>_internal.go`, pasta de comando, página em
`docs/Commands/`, exemplo em `examples/cli/<name>/`. Utilitários de edição em
`utils/command_edit.go`, par de `utils/route_edit.go` (`CommandClearSet`, `CommandFieldEdited`,
`CommandFieldEditEmpty`).

### 4.2 Alterados

| Comando | Mudança |
|---|---|
| `add-command` | `--help` e `--category` opcionais (defaults: `""` e `Commands`, como `add-route` usa `Routes`); escreve `command.yaml` + `InternalPureHandler.go`; aceita `--example`, `--long-description`, `--hidden` direto como `add-route` |
| `set-command` | `--clear <key>`; `--identifier` substitui a lista (hoje acrescenta) como `--method` em `set-route`; `--example` continua acrescentando, como em `set-route` |
| `add-flag` / `add-arg` | `--type` aceita `integer`/`number`; ganham `--enum` (repetível) e `--pattern` |
| `remove-command` | texto alinhado ao `remove-route`: "o build só renderiza; código à mão pode referir o que sumiu" |
| `help` | lê de `Entries`, como qualquer comando |

### 4.3 Nomes de props em `sandbox/api/`

| Rota | Comando |
|---|---|
| `RoutePathProps` / `RoutePathEditProps` | `CommandFieldProps` (era `FieldProps`) / `CommandFieldEditProps` |
| `RouteParameterEditProps` | — (flag e arg dividem o mesmo props, com `Kind`) |

## 5. Migração deste repo (~80 comandos)

O self-hosting exige que cada passo compile e seja idempotente.

| Passo | O que muda | Handlers antigos |
|---|---|---|
| 1 | `generated/cliio/`, `values.go` saindo do `climain.go`; `handle_*.go` da cli escritos uma vez; `api.CommandFailure`, `Cli.Fail` | intocados |
| 2 | `entries.go` gerado para todo comando (aditivo); `Items`/`Get*` continuam | intocados |
| 3 | `new.go` fecha sobre `InternalPureHandler` quando `verify` achar essa assinatura, senão sobre `CommandHandler` | coexistem |
| 4 | action oculta única `migrate-commands`: renomeia `entries.yaml`→`command.yaml`, `name:`→`id:`, `int/float`→`integer/number`, `handler.go`→`InternalPureHandler.go`, reescreve `command.GetX("id")`→`entries.Id` e o bloco de erro→`return err` | migrados |
| 5 | remove o caminho legado (`CommandHandler`, `Items`, `Get*`), remove `migrate-commands`; `check_commands.go` passa a exigir só a forma nova | — |
| 6 | verbos novos (§4.1) e alterados (§4.2) | — |
| 7 | `exec-test --update`; bump de versão | — |

O passo 3 é a única inferência pela forma de um arquivo e existe só entre os passos 3 e 5.

## 6. O que fica diferente, de propósito

| Mecânica da rota | Na cli | Por quê |
|---|---|---|
| `priority`, `--before/--after`, `rebalance-routes`, `phase` | não existe | argv casa exatamente um comando; não há cadeia |
| `trigger` / `--middleware` | não existe | idem; `--quiet` fica na dispatch (é o único "middleware" e é fixo) |
| `paths` (fatias do caminho) | `args` posicionais | já é o mesmo papel |
| `fonts` (query/header/cookie) | `identifiers` (`--out`, `-o`) | já é o mesmo papel |
| `body` + json-schema + `import-body` | não existe | argv não tem corpo; `--enum`/`--pattern` cobrem a validação útil |
| `response` / `response-type` | `Std.Printf` / `Std.Error` | canais já definidos no CLAUDE.md |

## 7. Arquivos a atualizar

| Arquivo | Mudança |
|---|---|
| `assets/sandbox-cli/**` | §3 |
| `assets/templates/command_*.{go,yaml,md}` | `command.yaml`, `InternalPureHandler.go`, `entries.go` (par de `route_entries.go`), `command_new.go` (par de `route_new.go`), `commandprops.go`, `cli_handle_*.go` |
| `sandbox/internal/actions/verify/check_commands.go` | novo: assinatura, nome do arquivo, chaves do yaml, ids únicos e exportáveis, tipos canônicos |
| `sandbox/internal/utils/command_{conf,edit,match}.go` | par de `route_{conf,edit,match}.go` |
| `assets/doc/docs/EntriesYaml/` → `CommandYaml/` em `assets/doc-cli/docs/` | par de `assets/doc-server/docs/RouteYaml/` |
| `assets/doc/docs/Rules/doc.md` | regras de handler (`error`, `cliio.Fail`, `handle_*` da cli) |
| `AgnosConfig/structure.yaml` | `sandbox/internal/cli/errors/`, `generated/cliio/`, `generated/cli/command/` |
| `docs/Contributing/doc.md` e `CLAUDE.md` (via template) | receita de comando nova, mesmo commit |
| `docs/Interview/doc.md` + interview | menus de `set-flag`, `rename-command`, `show-command` |
| `examples/cli/` | um exemplo por verbo novo; `exec-test --update` para os que mudam de forma |
