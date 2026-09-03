Crie o comando publish, ele ira olhar o project.yaml, entao buildar o projeto,compilar para todas as arquiteturas, entao usando o gh , ira criar uma release, subir os binarios.

comando:
```bash
 agnos publish --runtime go --publisher gh 
```


flags

--path, -p : The directory holding the project (defaults to the current directory)

--release-name, -rn : The name of the release

--draft : Create a draft release

--target, -t : The target to compile for (defaults to all)

--publisher, -pub : The publisher to use (defaults to gh)

## importante
se o user colocar um publisher que nao exista, o sistema deve falar qual publisher que pode ser usado
