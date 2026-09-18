## Mecanica de Database
Objetivo: Criar uma mecanica de geracao de banco de dados baseado em yaml
a ideia central e que a pasta sandbox/internal/databases contenha os bancos de dados da aplicacao, cada banco de dados tem que ser 100% gerado mecanicamente, a partir de um yaml em database/<banco>/specs.yaml , os methods tem que ser construidos baseados nas propriedades do banco.

## Referencias:
- fonte:SampleRepo/sandbox/internal/databases 
  - descricao: um exemplo base de como tem que ser cada banco de dados

### Importacoes:
 use o keep como depedencia (https://github.com/MateusMoutinhoOrg/Keep)

### Comandos
isso e so um guia, precisa reformular melhor
- database-init 
- database-purge
- database-add-table
- database-add-table-prop
- remove-table

