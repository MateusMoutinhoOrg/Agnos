# Relatório de Avaliação do Agnos CLI (Modo Entrevista)

## 1. Dificuldades com o Body JSON e Chaves (Schema)
O problema que você relatou sobre "não encontrar como determinar o body json e suas chaves" ocorre devido ao fluxo linear e fragmentado de como os comandos foram agrupados na interface interativa.

* **Descoberta do Comando (Discoverability):** O Agnos tem comandos dedicados a isso: `add-body-field` (para adicionar propriedades ao JSON) e `set-body` (para configurar Content-Type, obrigatoriedade, tamanho limite). No entanto, eles ficam isolados no menu **"Server System"**. Quando você cria uma rota (`add-route`), o CLI não engata perguntas sobre o body automaticamente; ele termina o comando e te joga de volta ao menu principal. Você precisa "adivinhar" que deve buscar por `add-body-field` depois.
* **Criação de Estruturas Complexas é Cansativa:** O `add-body-field` permite adicionar campos aninhados usando notação de ponto (ex: `address.city`), mas configurar um payload com 10 propriedades pelo modo entrevista exigiria rodar a opção "add-body-field" 10 vezes.
* **Falta de Suporte a Mock/Import:** **[Feature Faltante]** Seria muito útil ter a possibilidade de colar um JSON de exemplo e deixar o CLI inferir os campos (tipos, required, etc.), como um "import json".

## 2. Restrições Ocultas no Menu e Confusão (Ruled Out)
No código fonte (especificamente em `ruled_out.go`), o Agnos usa uma mecânica inteligente que esconde certas perguntas dependendo das respostas anteriores. 

* **Exemplo:** Se em `add-body-field` você responder que o tipo é `boolean`, perguntas sobre `min`, `max`, `pattern` ou `format` não serão exibidas.
* **O Problema:** Isso limpa o questionário, mas pode deixar a sensação de "cadê a feature de limites?" caso o usuário selecione acidentalmente um tipo incorreto. Ele só descobrirá na tela de confirmação (Confirm Screen). **[Sugestão de UX]** Seria bom um aviso sutil sobre os campos que foram inferidos ou ignorados por regra.

## 3. Gestão e Edição de Configurações (Falta de "Edit")
O modo entrevista lida muito bem com a "criação" (add-x) e "remoção" (remove-x), bem como purgar módulos inteiros. Mas falha no meio termo.

* **Falta de Edição Direta:** Não há opções fluidas como "edit-body-field". Se você configurou um campo `age` (idade) no JSON e esqueceu de colocar um limite `--max 130`, você precisará ir ao arquivo `route.yaml` manualmente, ou invocar o comando de remover (`remove-body-field`) e depois adicioná-lo de novo do zero na entrevista.

## 4. Retrato Visual (Preview) da Rota
Ao rodar vários comandos na mesma rota (ex: `add-param`, `add-header`, `add-body-field`), o usuário vai construindo o contrato da rota às cegas na entrevista.
* **[Feature Faltante]** Uma opção de "Visualizar Rota" (`view-route-schema`), onde o CLI imprima uma árvore (tree) de como a rota se encontra, com seus headers, parâmetros e árvore do body JSON.

## 5. Prós do Modo Entrevista
Apesar dos atritos nas configurações refinadas (como o JSON body), o Agnos Interview é tecnicamente muito bem feito nos seguintes aspectos:
* **Recuperação de Erros:** O loop de confirmação é brilhante (definido em `interview_internal.go`). Se o gerador rejeitar um campo inválido, ele não derruba toda a sua entrevista, ele permite alterar exatamente a resposta defeituosa (`PruneRuledOut`) e tentar rodar de novo.
* **Progressão Guiada:** A lógica do estado (`state.go`) de sugerir o próximo passo com a estrela `★` ("Declare its first route") impede que iniciantes se percam em categorias cujo módulo nem foi inicializado ainda.

## Resumo das Features Faltantes Recomendadas para o Roadmap
1. Sugerir configuração de body/headers imediatamente após rodar o `add-route`.
2. Criar um importador de Payload (Body JSON) baseado em um exemplo colado (inferência de tipo).
3. Opção de ver o "Resumo" da Rota (preview da árvore do schema configurado) por dentro do menu.
4. Adicionar comandos ou fluxo de edição para parâmetros de rota, query ou chaves JSON já declarados, diminuindo a dependência de `remove -> re-add`.
