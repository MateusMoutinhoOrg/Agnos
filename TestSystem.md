### Test Execution

Crie um sistema de exemplos que funcionará ao mesmo tempo de exemplo, mas também como teste.
Eu coloquei em `examples/` como a estrutura tem que funcionar.

### Comando e Ação a ser criado

~~~bash
agnos exec-test --path <path>
~~~

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
- Quando `<nome>` existe nos dois lados, os dois `result.yaml` devem ser idênticos: é a
  asserção de que a CLI é só um wrapper da lib.

### Funcionamento do comando

Ele tem que iterar sobre `examples/cli` e `examples/lib`, e para cada diretório executar a
seguinte ação:

- deletar a pasta `TestDir`
- se for lib, rodar o `example.go`
- se for cli, rodar o `example.sh`
- se houver o `result.yaml`
  - comparar os resultados do `result.yaml` com o resultado da execução
  - se forem iguais o teste passou, se não o teste reprovou
- se não houver o `result.yaml`
  - criar o `result.yaml` com os resultados que o teste gerou

Independente dos resultados, cada teste deve ser printado na cli, com o status
(a não ser que `--quiet` esteja presente).

O comando retorna `api.ExitOk` se todos os testes passarem e `api.ExitFailure` se qualquer um
reprovar.

### result.yaml

~~~yaml
exit-code: 0
cli-output: ""
tree:
 -
   file: AgnosConfig/docs/ReadmeHeader.md
   sha: "af0796675f5c16a2f2b5a3285d21406a30875803f12d51885c9cc33a882e6a51"
 -
   file: AgnosConfig/ignore.yaml
   sha: "63186ca147505db781c295629d94c763f1d0b5722b5854925af8f3a909c47255"
~~~

| campo | o que é |
|---|---|
| `exit-code` | status de saída do `example.sh` / `example.go`; `0` é sucesso |
| `cli-output` | stdout e stderr da execução, na ordem em que foram escritos |
| `tree` | todo arquivo dentro de `TestDir`, ordenado por `file`, com o sha256 do conteúdo |

Os três campos entram na comparação: uma execução que produz a mesma árvore mas sai com outro
`exit-code` reprova.

### O `agnos` que o exemplo chama

O `example.sh` escreve `agnos` porque é o que o usuário final digita — o exemplo tem que ler
como documentação, não como script de teste. Quem resolve isso é o `exec-test`: antes de rodar
a suíte ele compila `./cmd/main` para um binário temporário com o nome do projeto
(`name` de `AgnosConfig/project.yaml`) e coloca esse diretório na frente do `PATH` da execução,
via `rundeps.RunProps.Env`.

Consequências:

- o exemplo é testado contra o código atual do repo, nunca contra o binário instalado — é a
  mesma razão pela qual `build` é sempre bootstrapado;
- o mesmo alias serve o `example.go`, que importa a lib direto do módulo, então os dois lados
  exercitam a mesma árvore de fontes;
- em um projeto scaffoldado a regra é a mesma sem nenhuma adaptação: o alias tem o nome daquele
  projeto, apontando para o `cmd/main` daquele projeto.

### Documentação a ser criada

Duas docs, criadas com `agnos add-doc` e escritas como template (nunca à mão no `docs/`),
para renderizarem em todo projeto agnos:

| doc | template | temas | quando existe |
|---|---|---|---|
| `LibExamples` | `assets/all/docs/LibExamples/` | `reference`, `lib-usage` | sempre |
| `CliExamples` | `assets/cli/docs/CliExamples/` | `reference`, `cli-usage` | só se a cli estiver habilitada |

`CliExamples` fica no grupo `cli` justamente porque um projeto sem cli não tem `examples/cli/`
— o grupo já é renderizado só quando `sandbox/internal/cli` existe, e
`CollectGeneratedDocs` já indexa a doc a partir do `props.yaml` do grupo.

Cada uma documenta o seu lado do sistema, e só o seu lado:

- `LibExamples`: o layout de `examples/lib/<nome>/`, o contrato do `example.go`
  (`package main`, cwd = o próprio diretório, escreve só em `TestDir`), como adicionar um
  exemplo novo e como o `result.yaml` é gerado/comparado.
- `CliExamples`: o mesmo para `examples/cli/<nome>/` e `example.sh`, mais o alias `agnos` no
  `PATH` (o binário compilado de `./cmd/main`, com o `name` de `AgnosConfig/project.yaml`).

O que é comum aos dois — `agnos exec-test`, o schema do `result.yaml`, a regra de que
`<nome>` presente nos dois lados tem `result.yaml` idêntico — é dito uma vez só, e o outro
lado linka com caminho relativo.
