

##  priority mechanism 
O Mecanismo de de prioridade nas api.Route deve ser o seguinte: 
se houver as seguintes handlers com identifiers:
- /a/b/
- /a/b/c
- /
e supondo que a requisicao foi /a/b/c/d
nesse caso, deve se rodar a funcao handler do identifier "/", se ele nao fizer um write no body ou nao retornar um error, deve se rodar a funcao handler do identifier "/a/b/", se ele nao fizer um write no body ou nao retornar um error, deve se rodar a funcao handler do identifier "/a/b/c/"

