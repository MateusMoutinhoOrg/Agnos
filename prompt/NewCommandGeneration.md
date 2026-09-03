### Objetivo
refatorar o sistema de geracao de comandos , para ficar mais facil criar comandos

## Refs 
em prompt/sample_command voce pode ver como eu espero que cada pasta de comands seja, a ideia e que tenha 3 arquivos:
handler.go (codigo que o usuario colocara a logica do comand)
entries.go (struct criada  pelo agnos)
entries.yaml (arquivo que o usuario vai criar, contendo o funcionamento do comando)


## Remover 
### sandbox.Cli.Command
 (nao tem por que haver essa esstrutura mais)

## Fluxo de Build
ao rodar o comando agnos-build, o sistema tem que iterar sobre os arquivos de internal/commands, entao olhar o arquivo internal/commands/<item>/entries.yaml e para cada entries.yaml fazer a seguinte acao:
- criar o arquivo entries.go, que deve ser uma estrutura, com todos os elementos descritos no entries.yaml 


apos criar todos os entries.go, o build precisa criar o arquivo sandbox/internal/cli/climain.go , esse arquivo tem que validar todas as entradas (assim como esta agora), e caso seja o identifier, construir a estutura de entries, e chamar a funcao handler respectiva.


## Assets
Modifique os assets para ja ficarem nesse novo formato.

## Importante: 
nao refatore o proprio codigo em si, eu vou testar antes se esta funcionando, antes de aplicar o boostrap .

