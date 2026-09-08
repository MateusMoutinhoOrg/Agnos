### Refatoracao do Sistema de Teste



## Sistema de Update
O sistema de update, deve ser um comando especifico, chamado update-test,onde passa o caminho do teste, do jeito que ta, facilita muito regerar tudo de uma vez e nao ver onde estao os erros.

## Sistema de Coparation
por agora, todo o TestDir, e comparado, entao uma modificacao no start, vai quebrar o add-command por exemplo, o que e um erro, quero que adicioen a mecanica de comparation dir, apos gerar todos os resultados em TestDir, o example.sh. e o example.go, deve mover os arquivos relacionados ao teste para AssignatureDir, e a tree do result.yaml, tem que comparar o o ComparationDir, desse modo, somente arquivos que sejam ligados a aquele teste serao comparados, (diminuindo o acoplhamento entre os testes)