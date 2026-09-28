# Nova CLI

O `sandbox-cli` passa a espelhar o `sandbox-server`: um comando casa por **triggers** sobre o argv, e
a dispatch roda uma **cadeia** com middlewares.

## 1. Casamento

O argv se parte em **segmentos** (os tokens iniciais sem `-`, até o primeiro com `-`) e **flags**.
Triggers leem os segmentos; o verbo vem sempre primeiro.

| `command.yaml` | Par em `route.yaml` | Chaves |
|---|---|---|
| `args[]` | `paths[]` | `id`, `start`, `end`, `type` (`string`, `integer`, `number`, `uuid`; um intervalo liga `[]string`), `trigger`, `required`, `default`, `description` |
| `flags[]` | `parameters[]` | `id`, `keys` (as grafias que o usuário digita: `[--command, -c]`; default `[--<id em kebab-case>]`), `type` (`string`, `integer`, `number`, `boolean`, `string-array`, `integer-array`), `required`, `default`, `min`, `max`, `enum`, `pattern`, `trigger` |
| `priority`, `segments`, `category`, `help`, `long-description`, `examples`, `hidden` | idem | — |
| `strict` | — | default `true`; `false` em middleware |

- O texto de um trigger é o dos segmentos juntados por espaço. `prefix` é por segmento.
- Falhar um trigger ou um tipo em `args` é não-casamento. Faltar um `required` ou falhar a conversão de uma flag é erro de uso.
- `api.Trigger` e `MatchTrigger` passam a ser compartilhados pelas duas extensões (`assets/sandbox/sandbox/api/trigger.go`, `generated/trigger/`) e ganham o tipo `one-of` (`values: [...]`), usado quando um comando responde a mais de um nome.
- O default de `add-command <name>` é `args[0]` com `start: 0`, `end: 0`, `equal <name>`.

```yaml
priority: 100
args:
  - { id: Command, start: 0, end: 0, trigger: { type: equal, value: add-flag } }
  - { id: Name, start: 1, end: 1, required: true }
flags:
  - { id: Target, keys: [--command, -c], required: true }
category: Cli System
help: Add a flag to a command
```

`--pattern` compila como na rota:

| Peça | Vira |
|---|---|
| literais (`route add`) | um arg `equal` |
| `{name}` / `{name:integer\|number\|uuid}` | um segmento |
| `{*rest}` no fim | `start: i, end: -1` |
| sem `{*…}` | `segments: n` |

## 2. Cadeia

1. Coleta todo comando que casa.
2. Roda em `priority` crescente, com `Entries` ligado e `props` compartilhado.
3. O primeiro que responde encerra a cadeia.
4. Se ninguém respondeu, vai para `handle_not_found.go`.

Responder é `response.SetStatus(code)` ou `response.Printf(...)` (que fixa `ExitOk`).
`response.Error` e `response.Log` não respondem. Um handler que não responde recusou, e o próximo
roda.

**strict:** antes de um comando strict rodar, todo token que nem ele nem um middleware anterior
consumiu é erro de uso.

```go
func InternalPureHandler(sandbox *api.Sandbox, props *api.CommandProps, entries *Entries, response *api.CommandResponse) error
```

- Recusar é `return cliio.Fail(sandbox, api.ExitFailure, "<id>", "<msg>")`.
- `Entries` é gerado em `entries.go`, com um campo por arg e flag mais `FullCommand []string`.
- `props` vem de `sandbox/api/commandprops.go`, escrito uma vez.

## 3. Middlewares

`add-command --middleware` gera `prefix` vazio (casa todo argv), `priority: 10`, `strict: false` e um
stub que não responde.

| Gerado | Prioridade | Faz |
|---|---|---|
| `project` (no agnos) | 10 | declara `--path` e `--quiet`, põe `props.Path`, silencia `Log`; sai de todo comando |
| `help-flag` | 5 | `--help` imprime o help do próximo da cadeia; recusa se o próximo declara `--help` |
| `help`, `version` | 100 | comandos |
| `handle_not_found.go` | — | argv vazio → help geral |

## 4. Arquivos

```
sandbox/internal/commands/<snake>/
  command.yaml
  new.go                   gerado
  entries.go               gerado
  InternalPureHandler.go   escrito uma vez
```

| Novo | Par |
|---|---|
| `generated/cli/cli/{climain,new}.go` | `generated/server/server/` |
| `generated/cli/command/{new,CommandHandler,IsActionable}.go` | `generated/server/route/` |
| `generated/cliio/` | `generated/routeio/` |
| `sandbox/internal/cli/errors/handle_{not_found,bad_usage,unknown_flag,unexpected_arg,failure}.go` | `sandbox/internal/server/errors/` |
| `sandbox/api/commandprops.go`, `api.CommandResponse`, `api.CommandFailure`, `Cli.Fail` | `routeprops.go`, `RouteFailure`, `Server.Fail` |
| `utils/command_{conf,edit,match}.go`, `utils/trigger.go` | `utils/route_*.go` |
| `verify/check_commands.go` | `check_routes.go` |
| `assets/doc-cli/docs/CommandYaml/` | `RouteYaml` |
| `examples/cli/command-{chain,middleware,pattern}` | `route-*` |

## 5. Verbos

| Comando | Par | Flags |
|---|---|---|
| `add-command` | `add-route` | `--trigger*`, `--pattern`, `--middleware`, `--priority`, `--before`, `--after`, `--help`, `--category` |
| `set-command` | `set-route` | as mesmas + `--strict`, `--loose`, `--clear` |
| `remove-command`, `rename-command` | `remove-route`, `rename-route` | — |
| `rebalance-commands` | `rebalance-routes` | `--step` |
| `list-commands`, `show-command` | `list-routes`, `show-route` | só leem |
| `explain-command -- <argv…>` | `explain-route` | só lê |
| `add-arg`, `set-arg`, `remove-arg` | `*-path` | as de `add-path` + `--required`, `--default` |
| `add-flag`, `set-flag`, `remove-flag` | `*-parameter` | as de `add-parameter` sem `--font`, + `--key` (repetível), `--min`, `--max`, `--enum`, `--pattern` |

## 6. Ordem

1. `Trigger` compartilhado + `one-of`.
2. `cliio`, `handle_*`, `CommandResponse`, `commandprops.go`.
3. Dispatch em cadeia, `command.yaml`, `entries.go`, `InternalPureHandler`.
4. Migrar os comandos do agnos com uma action oculta `migrate-commands` e criar o middleware `project`; depois remover a action.
5. `check_commands.go`, verbos da §5, `help-flag`, `--pattern`.
6. Docs (`CommandYaml`, `Rules`, `Contributing`, `CLAUDE.md`: stdout do handler sai por `response`), interview, exemplos, `exec-test --update`, bump de versão.
