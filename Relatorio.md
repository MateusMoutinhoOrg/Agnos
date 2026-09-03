# Relatório de Bugs — Agnos 0.0.2

**Data:** 2026-09-03
**Binário testado:** `/usr/local/bin/agnos` (`agnos version` → `Version:0.0.2`)
**Ambiente:** macOS (darwin 22.6.0), Go 1.26.0
**Método:** exercício de todos os comandos (`start`, `build`, `verify`, `cli-init`, `cli-purge`,
`add-command`, `add-arg`, `add-flag`, `remove-*`, `set-command`, `dep-list`, `dep-install`,
`deps-init`, `help`, `version`) em projetos recém-criados, incluindo casos de erro e de borda.

---

## Resumo

| # | Severidade | Bug |
|---|-----------|-----|
| 1 | 🔴 Alta | `build` e `verify` reportam sucesso mesmo com projeto que não compila |
| 2 | 🔴 Alta | `add-command` aceita nomes inválidos e gera código Go inválido |
| 3 | 🔴 Alta | `--quiet` / `-q` não faz absolutamente nada em nenhum comando |
| 4 | 🔴 Alta | Flags desconhecidas são silenciosamente ignoradas (ex.: `--path` digitado errado) |
| 5 | 🟠 Média | `start` + `cli-init` deixam o projeto sem compilar (falta `go.sum` / `go mod tidy`) |
| 6 | 🟠 Média | Argumento `--array` descarta elementos inválidos silenciosamente (perda de dados) |
| 7 | 🟠 Média | `cli-purge` não remove `sandbox/internal/commands` e `sandbox/internal/cli` (deixa arquivo gerado órfão) |
| 8 | 🟠 Média | `add-arg --position` fora do intervalo / negativo dá mensagem de erro enganosa |
| 9 | 🟡 Baixa | Identificador gerado com erro de digitação: `CommandHander` |
| 10 | 🟡 Baixa | Códigos de saída inconsistentes entre erros de uso |
| 11 | 🟡 Baixa | Erros internos do Go vazam para o usuário final |
| 12 | 🟡 Baixa | Toda saída (inclusive erros) vai para stdout; stderr nunca é usado |
| 13 | 🟡 Baixa | Ruído de log misturado com saída "de máquina" (`dep-list`) |
| 14 | 🟡 Baixa | `start --module` documentado como opcional, mas é obrigatório sem `go.mod` |
| 15 | ⚪ Trivial | Falta de `\n` final na saída; `Agnos` vs `agnos`; nome de dep `serializebles` |

---

## Detalhes

### 1. 🔴 `build` / `verify` reportam sucesso com projeto quebrado

`agnos build` imprime `successfully rendered template` e `agnos verify` imprime `verify passed`
(exit 0 nos dois) mesmo quando o projeto **não compila**. A documentação do `build` diz
"compiling the source code into the output artifacts", mas nenhum dos dois invoca o compilador Go.

**Reprodução:**
```
agnos start -p demo -m x/demo
agnos add-command "Bad Name!" --help x --category Y     # ver bug #2
agnos verify        # -> "verify passed"  (exit 0)
agnos build         # -> "successfully rendered template"  (exit 0)
go build ./...      # -> malformed import path ".../bad_name!": invalid char '!'
```

**Esperado:** `build`/`verify` deveriam falhar (exit ≠ 0) quando o código gerado não é Go válido.

**Mudanças propostas:**

Para consertar o erro 1, adicione a flag `--runtime` — ex.: `--runtime "go"` (default `"go"`).
Então use o runtime `go` para executar comandos e garantir que o código está compilável.

---

### 2. 🔴 `add-command` aceita nomes inválidos e gera pacote Go inválido

Não há validação do nome do comando. Nomes são silenciosamente transformados
(`"Bad Name!"` → diretório `sandbox/internal/commands/bad_name!`, `MyCmd` → `mycmd`) sem aviso,
e caracteres ilegais são propagados para `package bad_name!` em `handler.go` / `entries.go`,
quebrando todo o build do projeto.

**Reprodução:**
```
agnos add-command "Bad Name!" --help x --category Y   # exit 0, cria pacote inválido
agnos add-command MyCmd --help x --category Y         # exit 0, vira "mycmd" sem avisar
```

**Esperado:** rejeitar nomes que não sejam identificadores válidos (`^[a-z][a-z0-9-]*$`),
ou pelo menos avisar sobre a normalização.

---

### 3. 🔴 `--quiet` / `-q` não funciona

Todos os comandos documentam `--quiet, -q` ("Quiets the cli output"). Em nenhum deles a flag
tem qualquer efeito — a saída é idêntica com e sem a flag.

**Reprodução:**
```
agnos start -p p -m x/p -q       # imprime "started with path .", "build started...", etc.
agnos build -q                    # imprime "verify started...", "verify passed", ...
agnos verify -q                   # imprime "verify started...", "verify passed"
agnos dep-install iodeps -q       # imprime tudo normalmente
```

---

### 4. 🔴 Flags desconhecidas são silenciosamente ignoradas

Qualquer flag não reconhecida é aceita e ignorada, com exit 0. Um erro de digitação em
`--path` faz o comando operar no diretório errado (o atual) sem qualquer aviso.

**Reprodução:**
```
agnos build --bogus-flag         # exit 0, roda normalmente
agnos build -z                    # exit 0
agnos verify --pathh /nao/existe  # exit 0 — verifica o diretório ATUAL, não /nao/existe
```

**Esperado:** `unknown flag "--pathh"` e exit ≠ 0.

---

### 5. 🟠 `start` + `cli-init` deixam o projeto sem compilar

Depois de `agnos start` seguido de `agnos cli-init` (ambos exit 0, "successfully rendered
template"), o projeto não compila porque o `go.sum` não é gerado/atualizado. É necessário
rodar `go mod tidy` manualmente.

**Reprodução:**
```
mkdir demo && cd demo
agnos start -p demo -m github.com/test/demo
agnos cli-init
go build ./...
# -> missing go.sum entry for module providing package
#    github.com/MateusMoutinhoOrg/Verb/sandbox
go mod tidy && go build ./...   # só agora funciona
```

**Esperado:** `agnos build` deveria rodar `go mod tidy` (ou `go mod download` + escrever
`go.sum`) para deixar o projeto em estado compilável.

---

### 6. 🟠 Argumento `--array` descarta elementos inválidos silenciosamente

Um argumento posicional `--array --type int` para de coletar no primeiro elemento inválido,
**sem erro** e descartando também os elementos válidos seguintes. O caso escalar equivalente
reporta o erro corretamente.

**Reprodução:**
```
agnos add-command calc --help c --category M
agnos add-flag n --command calc --type int
agnos add-arg items --command calc --array --type int
# ...
./demo calc --n 5 1 x 3     # -> items=[1]   (nenhum erro; "x" e "3" sumiram)
./demo calc --n abc 1        # -> erro correto: "abc" is not a valid integer
```

**Esperado:** erro `arg "items": "x" is not a valid integer` e exit ≠ 0.

---

### 7. 🟠 `cli-purge` deixa lixo para trás

Após `agnos cli-init` seguido de `agnos cli-purge`:

- `sandbox/internal/commands/version/entries.go` fica **órfão** — sem `entries.yaml` nem
  `handler.go` correspondentes (esses foram removidos, mas o arquivo gerado não).

A limpeza é parcial e inconsistente: mexe no comando `version` (criado pelo `start`, não pelo
`cli-init`) removendo uns arquivos e deixando outros.

**Esperado:** `cli-purge` deve remover por completo `sandbox/internal/commands` e
`sandbox/internal/cli`. As deps instaladas por `cli-init` (`sandbox/deps/argvdeps/`,
`sandbox/deps/std/`) **podem** ser deixadas para trás — não há problema nisso, já que podem
estar sendo usadas por outras partes do código.

---

### 8. 🟠 `add-arg --position` fora do intervalo dá erro enganoso

Quando existe um arg `--array` (que deve ficar por último), `--position` com valor grande
demais ou negativo (≠ -1) não reporta "posição inválida" — reporta a mensagem de
"array must stay last" com um índice sugerido incoerente.

**Reprodução:**
```
# comando calc já tem args: name (req), tags (--array)
agnos add-arg far --command calc --position 99   # -> "arg 'tags' ... insert before it with --position 2"
agnos add-arg neg --command calc --position -5   # -> mesma mensagem
```

**Esperado:** `--position 99 out of range (0..2)` / `--position must be >= 0`.

---

### 9. 🟡 Identificador gerado com erro de digitação: `CommandHander`

Todo `handler.go` gerado declara `func CommandHander(...)` (falta o "l" de "Handler"),
inclusive no comentário-doc. Compila, mas é um erro visível na API pública gerada.

```go
func CommandHander(deps *deps.Deps, entries *Entries) int {
```

---

### 10. 🟡 Códigos de saída inconsistentes

| Situação | Exit |
|----------|------|
| Erro de uso em `add-arg` / `add-flag` (tipo inválido, duplicado, etc.) | 2 |
| `start` sem flag obrigatória `--project-name` | 1 |
| Comando desconhecido (`agnos frobnicate`) | 1 |
| Flag desconhecida (`agnos build --bogus`) | 0 |

Convém padronizar (ex.: 2 para todo erro de uso/CLI).

---

### 11. 🟡 Erros internos do Go vazam para o usuário

- Fora de um projeto: `open go.mod: no such file or directory` (deveria ser algo como
  "no Agnos project found at <path>").
- Comando inexistente em `add-flag`: `command "nope" not found: open
  sandbox/internal/commands/nope/entries.yaml: no such file or directory`.
- Validação de int em runtime: `flag 'n': verb: "abc" is not a valid integer:
  strconv.Atoi: parsing "abc": invalid syntax` — vaza o nome do pacote interno (`verb:`)
  e o erro cru da stdlib.

---

### 12. 🟡 Tudo vai para stdout; stderr nunca é usado

Mensagens de progresso (`... started with path .`) **e** mensagens de erro
(`Unknown Command!`, `min must be an int`) são todas escritas em stdout. `stderr` fica sempre
vazio. Isso impede redirecionar log e resultado separadamente.

**Reprodução:**
```
agnos frobnicate 2>/dev/null    # ainda imprime "Unknown Command!"
agnos build 1>/dev/null          # não imprime nada — os logs foram para stdout
```

---

### 13. 🟡 Ruído de log misturado com saída de máquina

`agnos dep-list` imprime `dep-list started with path .` junto com a lista de deps, tudo em
stdout, dificultando o uso em scripts:
```
dep-list started with path .
embed
goimportsdeps
...
```
`agnos frobnicate` imprime só `Unknown Command!` (sem apontar `agnos help`), enquanto
`agnos help frobnicate` dá uma mensagem bem melhor (`✘ Unknown command: ... Run 'Agnos help'`).
Inconsistente.

---

### 14. 🟡 `start --module` documentado como opcional mas é obrigatório

`agnos help start` mostra `--module, -m` como `string │ optional`, mas sem `go.mod` no
diretório o comando falha:
```
agnos start -p demo
# -> the module flag (--module) is required when there is no go.mod in the path
```

---

### 15. ⚪ Diversos (trivial)

- Falta `\n` no fim de várias saídas — o prompt do shell cola no texto
  (`verify passedEXIT`, `successfully rendered template$`).
- Inconsistência de caixa: help e usage usam `Agnos` (maiúsculo) enquanto o binário é `agnos`.
- Nome de dep com erro de digitação em `dep-list`: `serializebles` (provável `serializables`).
- Struct `Entries` gerada ordena os campos em ordem alfabética (`Items` antes de `N`), não na
  ordem de declaração — surpreendente ao ler o código gerado.
- `go.mod` gerado fixa `go 1.25.0` mesmo com Go 1.26 instalado.

---

## O que funcionou bem

- Validação de `add-flag` / `add-arg`: tipo desconhecido, flag/arg duplicado, `--required`
  junto com `--default`, `--required` em boolean, `--min/--max` em string, `min > max`,
  `--min` não numérico — todos rejeitados com mensagens claras e exit 2.
- Validação de int em runtime para valor escalar (`flag 'n' must be <= 10`).
- `add-command` / `add-flag` / `add-arg` / `remove-command` recusam-se a sobrescrever/duplicar.
- Fluxo feliz (`start` → `cli-init` → `add-command` → `add-flag` → `add-arg` → `go mod tidy`
  → `go build` → executar) produz um CLI funcional.
- `--position 1` insere o arg na posição correta.
