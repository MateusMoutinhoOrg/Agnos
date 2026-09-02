## Group Render
create the group render mechanic 

### New Functions:

#### sandbox/internal/utils/render.go
~~~go 
func RenderGroup(deps *deps.Deps, io *smartio.SmartIO, path string, vars interface{}) error {

	return nil
}

~~~
a ideia da funcao rnder group, e renderizar um grupo de arquivos todos de uma vez, por exemplo chamando:

~~~go 
 RenderGroup(deps,io,"all",vars) 
~~~
vai renderizar todos os arquivos de assets/groups/all  e salvar em seus respectivos paths .

### Refatoracao de Actions:
refatore a action start para renderizar o grupo start. 
refatore o build para renderizar o grupo all , e o grupo deps (se a pasta sandbox/deps existir)

