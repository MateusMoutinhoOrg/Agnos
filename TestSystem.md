# Refatoracao do Sistema de Teste

Nota de trabalho, nao doc. Remover no commit que terminar a implementacao.

## Problema

`exec-test` grava em `result.yaml` a saida mesclada, o exit code e o sha256 de **todo**
arquivo do `TestDir` (`exec_tests_internal.go`, `treeOf`). Daí:

1. A tree do golden guarda o projeto inteiro (~35 arquivos por exemplo): uma mudanca em
   template de `start` quebra os 48 goldens de uma vez.
2. `--update` regenera todos os goldens juntos, entao o que quebrou nao aparece.

A saida ja e desacoplada (preambulo roda com `-q`; ver `examples/cli/add-command/example.sh`).
A tree tem que seguir a mesma ideia: **o exemplo declara o que e asserido**.

## 1. AssertDir

Nome: **`AssertDir`** — subconjunto do `TestDir` sobre o qual se assere. Nao `GoldenDir`:
o golden e o `result.yaml` (expectativa versionada), o `AssertDir` e resultado real,
git-ignored, regerado a cada run.

- O exemplo escreve so no `TestDir`; no fim **copia** para `AssertDir/` os arquivos ligados
  ao teste. Copia, nao move: o `TestDir` fica intacto para inspecao e nenhum passo seguinte
  perde arquivo.
- `exec-test` remove `AssertDir` antes de cada run, como ja faz com `TestDir`.
- A tree do `result.yaml` passa a ser a do `AssertDir`.
- `exec-test` **falha se o `AssertDir` estiver vazio** — senao esquecer uma copia faz o
  exemplo passar sem assertar nada.
- `assets/templates/example_cli.sh` e `assets/templates/example_lib.go` nascem com o par de
  copia escrito. Os dois lados copiam o mesmo conjunto, ou o cross-check cli-vs-lib quebra
  por motivo alheio ao codigo sob teste.
- `.gitignore`: `AssertDir` no lugar de `AssignatureDir` (orfao, ja adicionado).

Descartado: `exec-test` gravar so o que difere de um baseline — "qual e o baseline" fica
ambiguo quando o preambulo varia (`start` / `start + cli-init` / `start + deps-init`).

## 2. update-test

Comando `update-test <name>`, casca fina sobre a acao `exec_tests` com `update = true`.
A execucao continua em `ExecTestInternal(deps, path, only, update)`.

- `<name>` obrigatorio (sem nome = erro de uso). Nome + side, como `--only` /
  `add-cli-example` / `remove-cli-example`; nao caminho.
- Antes de sobrescrever, imprime o diff do golden: paths que entraram, sairam ou mudaram de
  sha, e saida antiga vs nova. Sem o diff o comando separado nao resolve nada.
- `exec-test --update` continua, para mudanca de shape em massa; deixa de ser o caminho normal.

## Propagacao

"sha256 de todo arquivo do TestDir" aparece em varios lugares e muda junto:

- `sandbox/internal/commands/help/handler.go` — long-descriptions de `exec-test`,
  `remove-cli-example`, `remove-lib-example`
- `sandbox/internal/commands/exec_test/entries.yaml`
- `assets/templates/example_cli.sh`, `assets/templates/example_lib.go`
- `assets/all/docs/` (o que for renderizado para todo projeto agnos)
- `sandbox/internal/utils/examples.go` — constante nova ao lado de `ExampleTestDir`
- `remove_cli_example` / `remove_lib_example` — apagar tambem o `AssertDir`

Ordem: `AssertDir` primeiro (muda todos os goldens de uma vez, com `exec-test --update`),
`update-test` depois.
