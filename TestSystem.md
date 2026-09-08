# Refatoracao do Sistema de Teste

Nota de trabalho, nao doc. Deve ser removida no commit que terminar a implementacao.

## Estado atual

`exec-test` roda cada `examples/{cli,lib}/<name>/`, remove o `TestDir`, executa
`example.sh` / `example.go` e grava em `result.yaml`: a saida mesclada, o exit code e o
sha256 de **todo** arquivo do `TestDir` (`exec_tests_internal.go`, `treeOf`). `--update`
reescreve todos os goldens de uma vez.

Dois problemas, independentes:

1. `--update` regenera tudo, entao o que quebrou nao aparece.
2. A tree do golden guarda o projeto inteiro (~35 arquivos por exemplo), entao uma
   mudanca em template de `start` quebra os 48 goldens de uma vez.

A saida ja e desacoplada: o preambulo roda com `-q`, entao `cli-output` so guarda a saida
do comando sob teste (ver `examples/cli/add-command/example.sh`). A tree tem que seguir a
mesma ideia: **o exemplo declara o que e asserido**.

## 1. update-test

Comando novo `update-test <name>`, casca fina sobre a mesma acao `exec_tests` com
`update = true`. Nao reimplementar a execucao: a logica continua em
`ExecTestInternal(deps, path, only, update)`.

- `<name>` e **obrigatorio** — sem nome e erro de uso. Nome + side, como
  `--only` / `add-cli-example` / `remove-cli-example`; nao caminho.
- Antes de sobrescrever, imprime o diff do golden: paths que entraram, sairam ou mudaram
  de sha, e a saida antiga vs a nova. E isso que faz ver onde estao os erros; separar o
  comando sozinho nao resolve nada.
- `exec-test --update` continua existindo para o caso de mudanca de shape em massa, mas
  deixa de ser o caminho normal.

## 2. AssertDir

Nome: **`AssertDir`** — e o subconjunto do `TestDir` sobre o qual se assere. Nao chamar de
`GoldenDir`: o golden e o `result.yaml` (a expectativa gravada, versionada); o `AssertDir` e
resultado real, git-ignored, regerado a cada run. Misturar os dois nomes apaga a distincao
de que o sistema inteiro depende. Nada de `AssignatureDir` / `ComparationDir`: nao sao
palavras em ingles e o doc anterior usava duas para a mesma coisa.

Mecanica:

- O exemplo continua escrevendo so no `TestDir`. No fim, `example.sh` / `example.go`
  **copia** para `AssertDir/` os arquivos ligados aquele teste. Copia, nao move: o
  `TestDir` fica intacto para inspecao quando um exemplo falha, e nenhum comando do
  exemplo corre o risco de perder um arquivo que um passo seguinte ainda usa.
- `exec-test` remove `AssertDir` antes de cada run, como ja faz com `TestDir` — senao
  sobra da run anterior passa como resultado da atual.
- A tree do `result.yaml` passa a ser a do `AssertDir`, nao a do `TestDir`.
- `exec-test` **falha se o `AssertDir` estiver vazio**. Sem isso, esquecer uma copia faz o
  exemplo passar sem assertar nada — o acoplamento cai junto com a cobertura, em silencio.
- Os templates `assets/templates/example_cli.sh` e `assets/templates/example_lib.go`
  nascem com o par de copia ja escrito. Os dois lados tem que copiar o mesmo conjunto, ou
  o cross-check cli-vs-lib quebra por um motivo que nao e o codigo sob teste.
- `AssertDir` entra no `.gitignore` no lugar de `AssignatureDir` (que ja foi adicionado e
  esta orfao). So os shas do `result.yaml` sao versionados.

Descartado: fazer o `exec-test` gravar sozinho so o que difere de um baseline. Nao exige
bookkeeping por exemplo, mas "qual e o baseline" fica ambiguo quando o preambulo varia
(`start` vs `start + cli-init` vs `start + deps-init`), e este repo prefere explicito a
magico.

## Propagacao

"sha256 de todo arquivo do TestDir" esta escrito em varios lugares e muda tudo junto:

- `sandbox/internal/commands/help/handler.go` — long-descriptions de `exec-test`,
  `remove-cli-example`, `remove-lib-example`
- `sandbox/internal/commands/exec_test/entries.yaml`
- `assets/templates/example_cli.sh`, `assets/templates/example_lib.go`
- `assets/all/docs/` (o que for renderizado para todo projeto agnos)
- `sandbox/internal/utils/examples.go` — constante nova ao lado de `ExampleTestDir`
- `remove_cli_example` / `remove_lib_example` — apagar tambem o `AssertDir`

Ordem: `AssertDir` primeiro (muda todos os goldens de uma vez, com
`exec-test --update`), `update-test` depois.
