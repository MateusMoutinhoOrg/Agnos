# Nova CLI — o sistema de comandos como espelho do sistema de rotas

Objetivo: o `sandbox-cli` passa a funcionar como o `sandbox-server`. Um comando deixa de ser
**achado por identificador** e passa a **casar por triggers** sobre as posições do argv. A dispatch
deixa de escolher um comando e passa a rodar uma **cadeia**, com middlewares, prioridade e fase
`after`. Onde o argv não tem par no http, o documento diz por quê (§8).

## 1. Diagnóstico — o que difere hoje

| Aspecto | Rota | Comando hoje |
|---|---|---|
| como casa | `paths[]`: fatias `start..end` dos segmentos, cada uma com `trigger` opcional; `parameters[]` com `trigger` | `identifiers: [add-flag]`: igualdade com `argv[0]` |
| quantos rodam | todas as rotas que casam, `priority` crescente, até uma responder | exatamente um |
| middleware | `--middleware`: rota que não responde e passa adiante | não existe; `--quiet` e `--help` estão fixos em `climain.go` |
| fase `after` | roda depois da resposta, lê `Entries.AnsweredStatus` | não existe |
| estado da cadeia | `props *api.RouteProps` (`sandbox/api/routeprops.go`, escrito uma vez) | nenhum |
| valores | `entries.go` gerado, `type Entries struct` com tag `id` | `command.GetString("path")` |
| handler | `InternalPureHandler(sandbox, props, entries, response) error` | `CommandHandler(sandbox, command) int` |
| falha | `routeio.Fail` → `routeio.Raise` → `sandbox/internal/server/errors/handle_*.go` do projeto | `Std.Error` + `ExitFailure` em cada handler; mensagens de uso fixas na dispatch |
| nada casou | `handle_not_found.go` do projeto | texto fixo em `climain.go` |
| `--pattern` | `/get-article/{article:integer}` compila em paths + `segments` | não existe |
| verbos | add/set/remove em cada seção, `rename-`, `list-`, `show-`, `explain-`, `rebalance-` | add/remove de flag e arg, `set-command` |
| `verify` | `check_routes.go` | sem `check_commands.go` |
| flags comuns | — | `--path` e `--quiet` repetidos no `entries.yaml` de todos os ~80 comandos |

## 2. O casamento — argv como caminho

### 2.1 Segmentos

O argv vira duas listas antes de qualquer comando ser consultado:

| Lista | Conteúdo | Par na rota |
|---|---|---|
| **segmentos** | os tokens que não começam com `-`, **até o primeiro que começa** | segmentos do caminho |
| **flags** | o resto, lido depois, por cada comando, pelas suas declarações | query / header / cookie |

`agnos add-flag output --command exec` → segmentos `[add-flag, output]`.

Um trigger só enxerga os segmentos iniciais. Isso resolve a ambiguidade `--out file` (valor da flag
ou posicional?) sem saber qual comando é: quem decide é a declaração de quem casou. Depois do
casamento, os `args` sem trigger (capturas) leem a lista posicional completa, com os mesmos índices,
já que os segmentos iniciais são um prefixo dela.

Custo: `agnos --path x build` deixa de casar `build`. O verbo vem sempre primeiro, e isso vira regra.

### 2.2 `args` = `paths`

| Chave | Efeito (idêntico a `paths[]`, salvo onde marcado) |
|---|---|
| `id` | campo de `Entries`, nome Go exportado, único entre `args` e `flags` |
| `start`, `end` | índices inclusivos, `-1` o último |
| `type` | `string`, `integer`, `number`, `uuid`; um intervalo (`start != end`) liga `[]string` **(difere: a rota junta com `/`)** |
| `trigger` | `{type, value, negate, ignore-case}`; texto = segmentos juntados por espaço **(difere: sem `/` inicial)** |
| `required`, `default` | **novo em relação a paths**: um arg sem trigger pode faltar. Faltar um `required` é erro de uso, e não não-casamento |
| `description`, `examples` | help |

`equal`, `prefix` (por segmento: `route` casa `route add`, nunca `router`), `text-prefix`, `suffix`
e `regex`, com a mesma semântica e o mesmo código de `api.Trigger`.

### 2.3 `flags` = `parameters`

| Chave | Efeito |
|---|---|
| `id` | campo de `Entries` |
| `key` | nome longo: `key: out` → `--out`. Default: `id` em kebab-case |
| `aliases` | grafias extras (`-o`) |
| `fonts` | ordem de leitura: `flag` (default), `env` (`<NAME>_<KEY>` em maiúsculas). É o par de `query`/`header`/`cookie` |
| `type` | `string`, `integer`, `number`, `boolean`, `string-array`, `integer-array`, os mesmos de `parameters` (o `--array` sai) |
| `required`, `default`, `min`, `max`, `enum`, `pattern` | `required` faltando ou conversão falha dá erro de uso (o `400` da cli) |
| `trigger` | a flag precisa casar para o comando rodar; falhar é não-casamento, como na rota |

### 2.4 Aliases de comando

`identifiers: [exec-test, test]` não cabe num `trigger.value` único. Proposta: um tipo de trigger
novo, `one-of`, com `values: [...]`, **em `api.Trigger`, compartilhado**. A rota ganha o mesmo
(`/users` ou `/people`). A alternativa sem tipo novo é `regex`, que fica ilegível no help.

## 3. `command.yaml` — mesmo formato de `route.yaml`

```yaml
priority: 100
phase: before            # omitido quando before
strict: true             # §4.3
segments: 2              # opcional, como na rota
args:
  - id: Command
    start: 0
    end: 0
    trigger: { type: equal, value: add-flag }
  - id: Name
    start: 1
    end: 1
    required: true
    description: the flag name
flags:
  - id: Target
    key: command
    aliases: [-c]
    required: true
  - id: Type
    key: type
    default: string
    enum: [string, integer, number, boolean, string-array, integer-array]
category: Cli System
help: Add a flag to a command's command.yaml
examples:
  - add-flag output --command exec
```

`add-command` sem `--trigger`/`--pattern` escreve `Command` (`start: 0`, `end: 0`, `equal <name>`).
**Difere da rota**, cujo default é o caminho inteiro (`start: 0`, `end: -1`): na cli os posicionais
seguintes são o caso comum.

| Hoje (`entries.yaml`) | Novo (`command.yaml`) |
|---|---|
| `identifiers` | `args[0]` com trigger (`one-of` se houver alias) |
| `name` em flag/arg | `id` + `key` |
| `identifiers` de uma flag | `key` + `aliases` |
| `type: int/float` + `array: true` | `integer`/`number`/`*-array` |
| `--path`, `--quiet` em todo comando | um middleware (§5), declarados uma vez |

O nome do arquivo muda para `command.yaml`: acaba a colisão entre `entries.yaml` (declaração) e
`entries.go` (struct gerado).

### 3.1 `--pattern`

| Peça | Vira |
|---|---|
| literais seguidos (`route add`) | um arg `equal`, nomeado por eles (`RouteAdd`) |
| `{name}` | um segmento, `Entries.Name string` |
| `{name:integer\|number\|uuid}` | um segmento tipado; não converter é não-casamento |
| `{*rest}`, só no fim | `start: i, end: -1`, `Entries.Rest []string` |
| sem `{*…}` | `segments: n` |

`agnos add-command route-add --pattern 'route add {name}'` dá **subcomandos de graça**, coisa que
`identifiers` não faz. Como na rota, o yaml nunca guarda o pattern; o `new.go` guarda o `Pattern`
derivado, que o `help` imprime como linha de uso.

## 4. A cadeia

### 4.1 Execução

Espelho de `servermain.go`:

1. Parte o argv em segmentos e flags (§2.1).
2. Coleta todo comando cujos `args` com trigger, `segments` e `flags` com trigger casam.
3. Roda em `priority` crescente. Cada um recebe `Entries` ligado e o `props` compartilhado.
4. O primeiro que **responde** encerra a cadeia.
5. Roda os de `phase: after` com o status congelado em `Entries.AnsweredStatus`.
6. Se ninguém respondeu, vai para `handle_not_found.go`.

### 4.2 Responder

Rota: `SetStatus` ou `Write` (que manda `200`). Cli: um `response *api.CommandResponse`, com

| Método | Efeito | Responde? |
|---|---|---|
| `SetStatus(code)` | código de saída | sim |
| `Printf(...)` | stdout; fixa `ExitOk` se nada foi fixado | sim |
| `Error(...)` | stderr | não |
| `Log(...)` | stderr, silenciado por `--quiet` | não |

Um handler que não faz nem `SetStatus` nem `Printf` **recusou**, e o próximo roda: é o que um
middleware é. O stub de `add-command` termina em `response.SetStatus(api.ExitOk)`. Um comando
silencioso (`build -q`) precisa dele, como a rota precisa do seu. **Muda a regra de canais do
CLAUDE.md:** dentro de um handler, stdout sai por `response`, e não por `Deps.Std.Printf`.

### 4.3 `strict`

Rota ignora query desconhecida; cli trata `--pathh` como erro. Solução: `strict: true` (default de
comando, `false` de middleware). Antes de um comando strict rodar, todo token não consumido por ele
**nem por um middleware anterior da cadeia** é erro de uso (`handle_unknown_flag` /
`handle_unexpected_arg`). Por isso flags globais funcionam: o middleware consome `--path`, e o
comando nunca o declara.

### 4.4 Assinaturas

```go
// InternalPureHandler.go — escrito uma vez por add-command, depois do projeto
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error {
	response.Printf("my-feature called\n")
	return nil
}
```

```go
// entries.go — gerado
type Entries struct {
	FullCommand []string `id:"FullCommand"` // o argv inteiro, sempre (par de FullRoute)
	Command     string   `id:"Command"`
	Name        string   `id:"Name"`
	Target      string   `id:"Target"`
}
```

- O handler retorna `error` e nunca um código. Recusar é `return cliio.Fail(sandbox, api.ExitFailure, "<id>", "<msg>")`. Some o `if err != nil { Std.Error; return ExitFailure }` repetido em ~80 handlers.
- `props *api.CommandProps` vem de `sandbox/api/commandprops.go`, que `build` escreve uma vez e o projeto tipa. Um middleware põe ali o que os de trás leem (`props.Path`).

## 5. Middlewares

`add-command --middleware` espelha `add-route --middleware`: `args` com trigger `prefix` vazio
(casa todo argv) salvo `--trigger`, `priority: 10`, `strict: false`, e um stub que não responde.

O que hoje está fixo e vira declaração:

| Hoje | Vira | Prioridade |
|---|---|---|
| `--path` e `--quiet` em ~80 `entries.yaml` | middleware `project` do agnos: declara as duas flags, valida o dir, `props.Path`, silencia `Log` | 10 |
| `quietId` fixo em `climain.go` | o mesmo middleware (a dispatch não conhece mais nenhuma flag) | — |
| `asksForHelp` / `runCommandHelp` | middleware gerado `help-flag`: flag `help` boolean com trigger `equal true`; imprime o help do próximo da cadeia; recusa se esse próximo declara `--help` (o caso de `add-command --help "…"`) | 5 |
| argv vazio → help | `handle_not_found.go`, escrito uma vez, imprime o help geral quando não há segmento | — |
| `help`, `version` | comandos gerados, como `health` | 100 |

Usos que o projeto ganha:

| Middleware | Trigger | Faz |
|---|---|---|
| grupo `route` | `prefix route` | carrega a config de rotas uma vez em `props` para `route add`, `route set`… |
| confirmação | `one-of remove-command remove-route …` | pergunta antes do destrutivo; responde `ExitFailure` se negado |
| env/auth | todo argv | lê token, recusa com `cliio.Fail` |
| `after` timing / log | todo argv, `phase: after` | registra duração e `AnsweredStatus` |
| aviso de versão | `phase: after` | "nova versão disponível" sem tocar no status |

O `help` lista, sob cada comando, as flags dos middlewares da frente dele: percorre a mesma cadeia
que `explain-command` percorre.

## 6. Camadas geradas

| Server | Cli nova |
|---|---|
| `generated/server/server/servermain.go` | `generated/cli/cli/climain.go`: só a cadeia (§4.1) |
| `generated/server/server/new.go` | `generated/cli/cli/new.go`: `Cli.Commands`, `Cli.Fail` |
| `generated/server/route/new.go` | `generated/cli/command/new.go`: base genérica |
| `generated/server/route/RequestHandler.go` | `generated/cli/command/CommandHandler.go`: argv → `Entries` via Reflectdeps |
| `generated/server/route/IsActionable.go` | `generated/cli/command/IsActionable.go`, lido por `utils/command_match.go` |
| `generated/routeio/` | `generated/cliio/` (`Fail`, `Raise`, `values.go`) |
| `sandbox/internal/server/errors/handle_*.go` | `sandbox/internal/cli/errors/handle_{not_found,bad_usage,unknown_flag,unexpected_arg,failure}.go`, escritos uma vez |
| `sandbox/api/routeprops.go` | `sandbox/api/commandprops.go`, escrito uma vez |
| `api.RouteFailure`, `Server.Fail` | `api.CommandFailure`, `Cli.Fail` |

**Compartilhado:** `api.Trigger` e `api.TriggerType` saem de `sandbox-server/.../api/route.go` para
`assets/sandbox/sandbox/api/trigger.go`, e `MatchTrigger` para `generated/trigger/`, renderizados
com qualquer uma das duas extensões ligada. `utils.RouteTrigger` e `RouteTriggerType` viram
`utils.Trigger` e `TriggerType`. Um só código de casamento para as duas.

## 7. Verbos — espelho um a um

| Rota | Comando | Flags |
|---|---|---|
| `add-route` | `add-command` | as mesmas: `--trigger*`, `--pattern`, `--middleware`, `--priority`, `--before`, `--after`, `--phase`, `--help`, `--category` (todas opcionais) |
| `set-route` | `set-command` | as mesmas + `--strict`/`--loose`, `--clear` |
| `remove-route` | `remove-command` | — |
| `rename-route` | `rename-command` | move o dir e reescreve `package` |
| `rebalance-routes` | `rebalance-commands` | `--step` |
| `list-routes` | `list-commands` | ordem da cadeia |
| `show-route` | `show-command` | árvore, não escreve |
| `explain-route <M> <path>` | `explain-command -- <argv…>` | quem roda, por que cada outro pula, qual `handle_*` responderia |
| `add-path` / `set-path` / `remove-path` | `add-arg` / `set-arg` / `remove-arg` | as de `add-path` + `--required`, `--default` |
| `add-parameter` / `set-parameter` / `remove-parameter` | `add-flag` / `set-flag` / `remove-flag` | as de `add-parameter` + `--alias`, `--min`, `--max`, `--enum`, `--pattern` |
| body, `import-body` | — | §8 |

Os nomes `arg`/`flag` ficam por serem o vocabulário da cli. As flags de cada par são as mesmas, e
os `utils/command_{conf,edit,match}.go` espelham `route_{conf,edit,match}.go`.

Sobre a armadilha do CLAUDE.md ("valor igual a uma grafia do próprio `add-flag`"): ela diminui,
porque `add-flag` perde `--identifier`. Mas continua em `--help`, agora dentro do middleware
`help-flag`.

## 8. O que fica diferente, de propósito

| Rota | Cli | Por quê |
|---|---|---|
| `methods` | — | argv não tem método |
| `body` + json-schema + `import-body` | — | argv não tem corpo; `enum`/`pattern` cobrem a validação útil |
| `response-type` | — | stdout não tem content-type |
| default de trigger no caminho inteiro | default em `start: 0, end: 0` | posicionais depois do verbo são o caso comum |
| intervalo liga texto com `/` | intervalo liga `[]string` | argv já vem separado |
| query desconhecida ignorada | `strict` | um typo de flag rodando no default é pior que um `404` |

## 9. Migração deste repo

| Passo | O que muda | Compatível com o antigo |
|---|---|---|
| 1 | `api.Trigger` compartilhado; o server passa a importar dele | sim |
| 2 | `generated/cliio/`, `handle_*.go` da cli, `CommandFailure`, `CommandResponse`, `commandprops.go` | sim |
| 3 | `command.yaml` lido quando existe, `entries.yaml` quando não; a dispatch vira cadeia; `identifiers` legado compila em `args[0]` `one-of` com `priority: 100`, `strict: true` | sim: um comando por argv, como hoje |
| 4 | action oculta `migrate-commands`: `entries.yaml`→`command.yaml`, remove `--path`/`--quiet` de cada um, cria o middleware `project`, `handler.go`→`InternalPureHandler.go`, `command.GetX("id")`→`entries.X`/`props.Path`, `Std.Printf`→`response.Printf`, bloco de erro→`return err` | — |
| 5 | remove o legado (`entries.yaml`, `Items`, `Get*`, `CommandHandler`) e `migrate-commands`; `check_commands.go` exige só a forma nova | — |
| 6 | verbos novos (§7), `help-flag`, `--pattern`, `one-of` | — |
| 7 | `exec-test --update`, bump de versão | — |

Os passos 3–5 leem dois formatos, e é a única inferência pela forma de um arquivo. Ela some no 5.

## 10. Arquivos a atualizar

| Arquivo | Mudança |
|---|---|
| `assets/sandbox-cli/**`, `assets/sandbox/sandbox/api/trigger.go` | §6 |
| `assets/sandbox-server/**` | importar o `Trigger` compartilhado; `one-of` |
| `assets/templates/command_*` | `command.yaml`, `InternalPureHandler.go`, `command_middleware_handler.go`, `command_entries.go`, `command_new.go`, `commandprops.go`, `cli_handle_*.go` |
| `sandbox/internal/actions/verify/check_commands.go` | novo, par de `check_routes.go` |
| `sandbox/internal/utils/command_{conf,edit,match}.go`, `trigger.go` | §7, §6 |
| `assets/doc/docs/EntriesYaml/` → `assets/doc-cli/docs/CommandYaml/` | par de `RouteYaml`, com seção "The chain" |
| `assets/doc/docs/Rules/doc.md` | cadeia, `response`, `strict`, canais |
| `AgnosConfig/structure.yaml` | dirs novos |
| `docs/Contributing/doc.md`, `CLAUDE.md` | receita nova e canais de saída, no mesmo commit |
| interview + `docs/Interview/doc.md` | `set-flag`, `rename-command`, middleware |
| `examples/cli/` | `command-chain`, `command-middleware`, `command-pattern`, par dos `route-*` |

## 11. Decisões em aberto

| # | Questão | Recomendação |
|---|---|---|
| 1 | `one-of` em `Trigger` ou `regex` para aliases | `one-of`, compartilhado |
| 2 | stdout via `response` ou manter `Deps.Std.Printf` e responder só por `SetStatus` | `response`, pela paridade com `Write` |
| 3 | renomear `add-arg`/`add-flag` para `add-path`/`add-parameter --command` | não, o vocabulário da cli fica |
| 4 | verbo sempre primeiro (§2.1) | aceitar e virar regra |
