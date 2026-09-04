### Test Execution 
Crie um sistema de exemplos que  funcionara ao mesmo tempo de exemplo , mas tambem como teste.
eu coloquei em examples , como a estrutura tem que funcionar.

### Comando e Acao a ser criado

~~~bash
agnos exec-test --path <path> 
~~~

### Funcionamento do comando
ele tem que iterar sobre examples/cli e examples/lib, e para cada diretorio executar a seguinte acao.
- deletar a pasta TestDir 
- se for lib, roadar example.go
- se for cli rodar o example.sh 
- se houver o result.yaml 
  - comparar os resultados do result.yaml , com o resultado da execucao,
  - se forem iguais o teste passou, se nao o teste reprovou 
- se nao houver o result.yaml
  - criar o result.yaml com os resultados que o teste gerou 


independente dos resultados, cada teste deve ser printado na cli, com o status 
(a nao ser que --quiet esteja presente)
