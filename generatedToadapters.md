
## Problem:
sandbox tem muitos intens e funcoes que nao estao ligados a regra de negocio do usuario, o que dificulta muito para llms e pessoas entenderem o codigo.

## A fazer
mapear no sandbox todos os itens que nao se encaixam nas regras de sandbox, e transformar eles em libs opinated. 

## Regra do sandbox:
no sandbox deve existir somente coisas relacionadas a regras de negocio do projeto, como rotas, implementacao de comandos, schema de banco de dados, etc.

## O que mudar:
(praticamente tudo de internal/generated) sao funcoes utils, que poderiam ser portadas para libs.


## OpinatedLibs 
como essas libs sao altamente opinativas, elas devem seguir esse padrao:
OpinatedAgnosCli , OpinatedAgnosServer, etc... 

## Importante: 
nao migre coisas que o user possa querer editar, como schema de banco de dados, rotas, etc...