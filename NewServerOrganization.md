Usand Agnosteset como base do esqueleto, restruture a mecanica de servers

## arquivos:
- sandbox/internal/routelist/<route>/new.go
esse construtor tem que ser uma imagem 1:1 do  route.yaml
- sandbox/internal/routelist/<route>/entries.go 
esse arquivo gerado, deve ser uma struct criada a partir das configuracos da rota

- sandbox/internal/routelist/<route>/InternalPureHandler.go
deve receber entries como argumento

### Importante
o parametro de priority, assim como o response type devem ser itens obrigatorios de todos os routes.yaml 
