### Test Execution

Crie um sistema de exemplos que funcionará ao mesmo tempo de exemplo e de teste.
`examples/` já tem a estrutura de referência.

### Comando e ação a serem criados

~~~bash
agnos exec-test --path <path> [--only <nome>] [--update] [--quiet]
~~~

| camada | arquivo | papel |
|---|---|---|
| ação | `sandbox/internal/actions/exec_test/exec_test.go` | valida props, loga, chama o internal |
| ação | `sandbox/internal/actions/exec_test/exec_test_internal.go` | lógica pura sobre `path` |
| api | `sandbox/api/actions.go` | `ExecTest func(props ExecTestProps) error` + `ExecTestProps{Path, Only, Update string/bool}` |
| comando | `sandbox/internal/commands/exec_test/` | `entries.yaml` (via `add-command`/`add-flag`), `entries.go`, `handler.go` |

O molde é `compile`: recebe `path` cru, junta caminhos relativos ao projeto como string, e usa
`deps.Rundeps`/`deps.Iodeps` direto. **`exec-test` não abre SmartIO** — não há transação a
persistir, o `TestDir` é escrito por um processo filho (fora do buffer) e o `tree` precisa da
árvore literal, sem os filtros de `ignore.yaml`/`paths.yaml`. Consequência: o `path` do projeto
é juntado na própria ação, não no boundary do SmartIO.

`handler.go` devolve `api.ExitOk` quando todo teste passa e `api.ExitFailure` quando qualquer um
reprova (a ação devolve `error` com o resumo; o handler imprime e converte, como em `compile`).

### Estrutura dos examples

~~~
examples/
  cli/                        # um diretório por comando exercitado pela CLI
    start/
      example.sh              # obrigatório: o exemplo, roda com cwd = este diretório
      result.yaml             # golden; criado no primeiro run se não existir
      TestDir/                # gerado pela execução, gitignorado, apagado antes de cada run
    add-command/
      example.sh
      result.yaml
  lib/                        # um diretório por ação exercitada pela lib
    start/
      example.go              # obrigatório: `package main`, roda com cwd = este diretório
      result.yaml
      TestDir/
    add-command/
      example.go
      result.yaml
~~~

Regras da estrutura:

- `examples/cli/<nome>/` contém exatamente um `example.sh`; `examples/lib/<nome>/` contém
  exatamente um `example.go`.
- Todo exemplo escreve apenas dentro do próprio diretório, em `TestDir`.
- Quando `<nome>` existe nos dois lados, `tree` e `exit-code` dos dois `result.yaml` devem ser
  idênticos: é a asserção de que a CLI é só um wrapper da lib. `cli-output` **não** entra nessa
  comparação cruzada — cada lado tem o seu texto (erro da CLI de um lado, `panic` do outro) e é
  comparado só contra o próprio golden.
- Um exemplo nunca é criado nem apagado à mão: `add-cli-example` / `add-lib-example` /
  `remove-cli-example` / `remove-lib-example` (abaixo), mesma regra de `add-doc`/`remove-doc`.
- `examples` entra em `AgnosConfig/structure.yaml` (a árvore de `docs/Structure` é gerada de lá);
  `TestDir` já está no `.gitignore`.

### Funcionamento do comando

Para cada `<nome>` de `examples/cli` e `examples/lib`, em ordem alfabética, lado `cli` antes de
`lib`:

1. `deps.Iodeps.RemoveDir(<dir>/TestDir)` — remoção direta em disco, sem buffer.
2. rodar o exemplo com `deps.Rundeps.Run`, `Dir = <dir>` absoluto/relativo ao `--path`:
   `sh example.sh` no lado cli, `go run example.go` no lado lib.
3. montar o resultado da execução: `exit-code` e `cli-output` vêm do `rundeps.Result`;
   `tree` vem de `ListFilesRecursively(<dir>/TestDir)`, ordenado por `file`.
4. se existe `result.yaml` e `--update` não foi passado: comparar; senão, gravar o
   `result.yaml` com o resultado gerado.

Cada teste é impresso com o status via `deps.Std.Log` (silenciado por `--quiet`). Uma reprovação
imprime **o que divergiu**: `exit-code` esperado/obtido, diff de `cli-output`, e por arquivo do
`tree` a marca `+` (só no obtido), `-` (só no golden) ou `~` (sha diferente) — comparar duas
árvores de centenas de arquivos e dizer só "reprovou" é inútil.

`--only <nome>` roda um único `<nome>` (os dois lados). `--update` regrava os goldens de tudo que
rodou: sem ele, o único jeito de atualizar um golden é apagar o `result.yaml` na mão.

### result.yaml

~~~yaml
cli-output: "started with path TestDir \n"
exit-code: 0
tree:
    - file: AgnosConfig/docs/ReadmeHeader.md
      sha: "af0796675f5c16a2f2b5a3285d21406a30875803f12d51885c9cc33a882e6a51"
    - file: AgnosConfig/ignore.yaml
      sha: "63186ca147505db781c295629d94c763f1d0b5722b5854925af8f3a909c47255"
~~~

| campo | o que é |
|---|---|
| `cli-output` | stdout e stderr da execução, na ordem em que foram escritos, normalizado |
| `exit-code` | status de saída do `example.sh` / `example.go`; `0` é sucesso |
| `tree` | todo arquivo dentro de `TestDir`, ordenado por `file`, com o sha256 do conteúdo |

Os três campos entram na comparação contra o próprio golden: uma execução que produz a mesma
árvore mas sai com outro `exit-code` reprova.

O arquivo é escrito por `deps.Serializables.SerializeToYaml`, com chaves em ordem alfabética e
sem espaço à direita — o golden é regenerado, então tem que sair byte-a-byte igual ou a suíte
deixa de ser idempotente. Os `result.yaml` hoje em `examples/` são placeholders escritos à mão:
estão sem `exit-code` e com `cli-output: ""`, que não bate com o `Std.Log` que `start` emite.

**Normalização do `cli-output`** (sem ela o golden reprova em outra máquina): antes de comparar
ou gravar, substituir toda ocorrência do caminho absoluto do diretório do exemplo por `<dir>` e
normalizar `\r\n` para `\n`. Um exemplo cuja saída carrega qualquer outro caminho absoluto, tempo
ou versão resolvida não é um exemplo válido.

**Volatilidade do `tree`**: um exemplo que dispara o runtime go (`start` roda `build`, que roda
`go mod tidy` e `go build`) produz `go.sum` dependente das versões que o proxy resolve naquele
dia. Ou o exemplo evita o runtime, ou `go.sum` e `release/` ficam fora do `tree`. Decidir isso
faz parte da implementação, não depois.

### sha256

Não existe capability de hash nas `deps`. `crypto/sha256` é stdlib pura, não é OS-bound, e o
sandbox já importa stdlib pura em `internal/` (`go/ast`, `go/parser`, `sort`): a ação importa
`crypto/sha256` direto. A alternativa — um `sandbox/deps/hashdeps` com adapter e `Bind` — só se
justifica se a regra passar a ser "nenhum import stdlib fora dos já usados".

### O `agnos` que o exemplo chama

O `example.sh` escreve `agnos` porque é o que o usuário final digita — o exemplo tem que ler como
documentação, não como script de teste. Quem resolve isso é o `exec-test`: antes de rodar a suíte
ele escreve, num diretório temporário, um script executável com o `name` de
`AgnosConfig/project.yaml` (aqui, `agnos`):

~~~sh
#!/bin/sh
exec go run <path-absoluto>/cmd/main "$@"
~~~

e põe esse diretório na **frente** do PATH da execução. Um alias `go run` em vez de um binário
compilado dispensa a etapa de compilar e o cleanup do binário; o cache de build do go faz a
segunda invocação em diante custar quase nada. O preço: `go run` escreve `exit status N` no
stderr quando o programa sai não-zero, então essa linha faz parte do `cli-output` dos exemplos
que falham de propósito. Compilar uma vez para o temp dir é a alternativa, e troca essa linha
por uma etapa de build no começo da suíte.

**Mudança de contrato necessária.** Pôr o diretório na frente do PATH não é possível com
`rundeps.RunProps.Env` como ele existe hoje: `adapters/libs/rundeps/rundeps.go` faz
`append(os.Environ(), props.Env...)` e o `os/exec` mantém a última ocorrência de uma chave
duplicada, então mandar `PATH=<tmp>` **substitui** o PATH inteiro — o `example.sh` perde `sh` e
`go`, e o lado lib nem inicia. O sandbox também não pode ler o PATH atual (sem `os`). A correção
é um campo novo:

~~~go
// PathPrefix are directories prepended to the PATH the program sees, ahead
// of the inherited one. The adapter is what reads the current PATH and joins
// it: the sandbox cannot.
PathPrefix []string
~~~

resolvido em `adapters/libs/rundeps/rundeps.go`. Mantém o sandbox puro, e é preferível a expor um
`Getenv` nas deps.

O `name` do projeto é usado literalmente como nome do alias, e o exemplo tem que digitar
exatamente esse `name`: num projeto com nome capitalizado, um `example.sh` que digita minúsculo
passa no macOS e reprova no Linux.

Consequências:

- o exemplo é testado contra o código atual do repo, nunca contra o binário instalado — é a mesma
  razão pela qual `build` é sempre bootstrapado;
- o mesmo alias serve o `example.go`, que importa a lib direto do módulo, então os dois lados
  exercitam a mesma árvore de fontes;
- em um projeto scaffoldado a regra é a mesma sem nenhuma adaptação: o alias tem o nome daquele
  projeto, apontando para o `cmd/main` daquele projeto.

### Comandos de gerenciamento

~~~bash
agnos add-cli-example <nome> --path <path>
agnos remove-cli-example <nome> --path <path>
agnos add-lib-example <nome> --path <path>
agnos remove-lib-example <nome> --path <path>
~~~

O molde é `add-doc`/`remove-doc`, não `exec-test`: estes escrevem na árvore do projeto, então
**abrem SmartIO**, chamam o `*Internal`, dão `Persist` e rodam `build` como follow-up com
`api.RuntimeNone` — `examples/` não está em `compilableDirs`, nada a compilar, mas o listing de
exemplos das docs precisa ser reescrito (ver abaixo).

| ação | escreve | recusa quando |
|---|---|---|
| `add_cli_example` | `examples/cli/<nome>/example.sh` | `<nome>` já existe; o projeto não tem `sandbox/internal/cli` |
| `add_lib_example` | `examples/lib/<nome>/example.go` | `<nome>` já existe |
| `remove_cli_example` | apaga `examples/cli/<nome>/` inteiro | `<nome>` não existe |
| `remove_lib_example` | apaga `examples/lib/<nome>/` inteiro | `<nome>` não existe |

Os `remove-*` apagam o diretório inteiro — `result.yaml` e `TestDir` junto. Um projeto sem cli
não tem `examples/cli/`, então `add-cli-example` falha com erro de uso em vez de criar o
diretório.

Cada um é o par de sempre: `sandbox/internal/actions/<nome>/<nome>.go` + `<nome>_internal.go`,
`sandbox/internal/commands/<nome>/` declarado com `agnos add-command` + `agnos add-flag` (nunca
editando o `entries.yaml` à mão), e a entrada em `sandbox/api/actions.go`:

~~~go
AddCliExample    func(path string, name string) error
RemoveCliExample func(path string, name string) error
AddLibExample    func(path string, name string) error
RemoveLibExample func(path string, name string) error
~~~

Argumento `<nome>` obrigatório; flags `--path` (default `.`) e `--quiet`, como todo comando.
Categoria `Examples` no `entries.yaml`.

**Templates do scaffold**, single-file em `assets/templates/`, na convenção de
`templates/doc_doc.md` e `templates/command_handler.go`:

| template | renderiza | vars |
|---|---|---|
| `templates/example_cli.sh` | `examples/cli/<nome>/example.sh` | `Name`, `ProjectName` (o alias que o script digita) |
| `templates/example_lib.go` | `examples/lib/<nome>/example.go` | `Name`, `Module` (o import path do projeto) |

O stub gerado tem que rodar e sair com `0` sem edição: o primeiro `exec-test` depois de um
`add-*-example` não pode reprovar, ele grava o golden.

**Listing gerado nas docs.** `sandbox/internal/utils/examples.go` expõe um `CollectExamples` que
lê `examples/cli` e `examples/lib` e devolve os `<nome>` ordenados, do mesmo jeito que
`CollectGeneratedDocs` lê as docs. `LibExamples` e `CliExamples` renderizam a tabela de exemplos
existentes a partir dele — documentar por listagem gerada, não escrevendo o nome de cada exemplo
à mão. É por isso que os quatro comandos rodam `build` no fim.

### Documentação a ser criada

Duas docs, criadas com `agnos add-doc` e escritas como template (nunca à mão no `docs/`), para
renderizarem em todo projeto agnos:

| doc | template | temas | quando existe |
|---|---|---|---|
| `LibExamples` | `assets/all/docs/LibExamples/` | `reference`, `lib-usage` | sempre |
| `CliExamples` | `assets/cli/docs/CliExamples/` | `reference`, `cli-usage` | só se a cli estiver habilitada |

`CliExamples` fica no grupo `cli` justamente porque um projeto sem cli não tem `examples/cli/` —
o grupo já é renderizado só quando `sandbox/internal/cli` existe, e `CollectGeneratedDocs` já
indexa a doc a partir do `props.yaml` do grupo (`assets/cli/docs/CliInstall` é o precedente).

Cada uma documenta o seu lado do sistema, e só o seu lado:

- `LibExamples`: o layout de `examples/lib/<nome>/`, o contrato do `example.go` (`package main`,
  cwd = o próprio diretório, escreve só em `TestDir`), `agnos add-lib-example` /
  `remove-lib-example` como o único jeito de criar e apagar um exemplo, e como o `result.yaml` é
  gerado/comparado.
- `CliExamples`: o mesmo para `examples/cli/<nome>/`, `example.sh` e
  `add-cli-example`/`remove-cli-example`, mais o alias `agnos` no PATH (`go run ./cmd/main`, com
  o `name` de `AgnosConfig/project.yaml`).

O que é comum aos dois — `agnos exec-test`, o schema do `result.yaml`, a normalização do
`cli-output`, a regra de que `<nome>` presente nos dois lados tem `tree` e `exit-code` idênticos —
é dito uma vez só, em `LibExamples`, e `CliExamples` linka com caminho relativo.

Num projeto scaffoldado o `examples/` nasce do primeiro `add-lib-example` / `add-cli-example`:
`start` não semeia nada e `assets/start/` não ganha diretório novo. As docs descrevem um
diretório que pode ainda não existir, e dizem qual comando o cria — o mesmo que `docs/` faz com
`add-doc`.
