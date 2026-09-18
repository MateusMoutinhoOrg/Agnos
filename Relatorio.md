# Relatório de Bug: Falso Positivo de "Drift" no Agnos CLI (`verify`)

## 1. Informações do Ambiente
- **Ferramenta:** Agnos CLI
- **Versão:** `v0.8.1`
- **Comandos afetados:** `agnos verify`, `agnos add-dep`, `agnos set-dep`

## 2. Descrição do Problema
O comando `agnos verify` reporta um erro de "drift" (desvio) em dependências remotas mesmo quando elas acabaram de ser baixadas com sucesso pela própria CLI, gerando um falso positivo. Isso afeta o fluxo de verificação do projeto e cria um cenário insolúvel ("Catch-22").

O bug ocorre especificamente com dependências remotas importadas com uso de renomeação de pacote (flag `--as`), onde o gerador interno de código da CLI altera legitimamente o `sandbox.go` para remover o field `Deps` e sua importação (para manter o sandbox isolado). 

O `agnos verify` não reconhece que essa alteração é um comportamento intencional (by design) e a aponta como uma diferença em relação ao módulo original remoto.

## 3. Passos para Reproduzir (Reproducer)
1. Baixe uma dependência remota utilizando a flag `--as` para renomear o diretório/pacote. 
   **Exemplo:**
   ```bash
   agnos add-dep github.com/MateusMoutinhoOrg/Keep --as database
   ```
2. Após o download concluir, rode o comando de verificação:
   ```bash
   agnos verify
   ```
3. A CLI retornará a seguinte violação na saída:
   ```text
   verify found 1 violation(s):
     - sandbox/deps/database/sandbox.go has drifted from the module it was copied from (run `agnos set-dep database --version <version>`)
   ```
4. Se tentarmos consertar executando `agnos set-dep database --version v0.0.7`, a mesma mensagem continuará aparecendo em testes subsequentes do `verify`.

## 4. Análise Técnica (O Efeito "Catch-22")
A CLI do Agnos tem a regra estrita de que **nenhum contrato dentro de `sandbox/deps/` pode realizar importações de outros pacotes** ou trazer fiação extra (wiring).
Portanto, o comando de geração (`add-dep` / `set-dep`) exclui as seguintes linhas do arquivo `sandbox.go` que vieram da dependência remota:
```go
import (
	"github.com/MateusMoutinhoOrg/Keep/sandbox/deps"
)
```
e o campo respectivo dentro da struct `Sandbox`:
```go
	Deps      *deps.Deps
```

Ao inspecionar as diferenças dos arquivos (diff) antes e depois da instalação da dependência, percebe-se um **conflito interno na CLI**:
- **Se o arquivo for idêntico ao remoto** (com a importação e o `Deps` presentes), o `verify` trava reclamando de importações não permitidas:
  ```text
  sandbox/deps/database/sandbox.go imports github.com/MateusMoutinhoOrg/Keep/sandbox/deps; sandbox/deps/<x>/ may import nothing at all
  ```
- **Se o arquivo sofreu a filtragem do gerador** (com a importação e o `Deps` apagados), o `verify` trava reclamando que o arquivo "desviou" (drifted) do módulo original na web.

## 5. Conclusão
O comparador de drift do `agnos verify` (que analisa AST ou diff) falha em ignorar intencionalmente a ausência do campo `Deps` e sua respectiva importação para dependências instaladas através da feature de "aliasing" (renomeação remota com a flag `--as`). Como resultado, qualquer dependência aliased ficará eternamente presa nesse erro do validador.
