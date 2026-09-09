## Referencia base
front-model 

### Comandos Novos: 

### front-init

#### exemplo: 
~~~bash
agnos front-init
~~~

#### Fluxo:
1. iniciar a estrutura esqueletica de assets eg.(front-model/assets/frontend)
2. Copiar as funcoes de utils de front-model ref(sandbox/internal/utils/templates.go)
3. Copiar  a rota static de  front-model

### add-page
#### Definicao:
adiciona uma nova pagina a aplicacao
#### exemplo: 
~~~bash
agnos add-page home
~~~
#### Fluxo:
1. Criar a rota /<pagina> semelhante a (front-model/sandbox/internal/routes/home)
2. criar um html modelo em assets/frontend/pages/<pagina>.html

### remove-page
#### Definicao:
remove uma pagina da aplicacao
#### exemplo: 
~~~bash
agnos remove-page home
~~~
#### Fluxo:
1. Remover a rota /<pagina> de (front-model/sandbox/internal/routes/home)
2. remover o html modelo em assets/frontend/pages/<pagina>.html