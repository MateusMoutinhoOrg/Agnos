### Objetivo
refatorar o sistema de geracao de comandos , para ficar mais facil criar comandos

## Remover 
### sandbox.Cli.Command
 (nao tem por que haver essa esstrutura mais)

## Fluxo de Build
ao rodar o comando agnos-build, o sistema tem que iterar sobre os arquivos de internal/commands, entao olhar o arquivo internal/commands/<item>/entries.yaml e para cada entries.yaml fazer a seguinte acao:
- criar o arquivo entries.go, que deve ser uma estrutura, com todos os elementos descritos no entries.yaml 


apos criar todos os entries.go, o build precisa criar 
