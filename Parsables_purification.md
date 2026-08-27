
## Purificacao dos parsables

Itere sobre os arquivos de sandbox/parsables e faca as seguintes alteracoes neles.

1. todos os modulos devem receber o uma string como entrada, ao ives do path, esssa string tem que chamar content.


2. remova todos os metodos de persist.

3. adicione o metodo render, que ira retornar uma string com o conteudo renderizado.


## Ideia
a ideia e quue todos os parsables sejam puros, eles recebem uma string inicial, parseaiam ela, e tem um metodo que renderiza de volta, assim eles nao fazem nenhum tipo de leitura e escrita ,


### Importante:
todos os elementos devem ter dois **New** um aceitando uma string de content, e outro vazio, assim :

~~~go 

func NewThemesConf(sandbox *api.SandBox, path string) (*ThemesConf, error)

func NewThemesConfEmpty(sandbox *api.SandBox) *ThemesConf 

~~~