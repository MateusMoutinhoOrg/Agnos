

### Embed mechanic 

### Objective:
setup a embed mechanic were sandbox will be able to retrive embed files,from a given directory 



### Ideia:
The ideia of embed assets its to store text,images,and templates for itens.



###  File Modifications:
- sandbox/contracts/deps/embed/embed.go
  - [example](./tree/sandbox/contracts/deps/embed/embed.go/example.go)

- sandbox/contracts/deps/deps.go
  - [modifications](./tree/sandbox/contracts/deps/deps/modifications.md)

- adapters/standard/embed.go
  - [example](./tree/adapters/standard/embed.go/example.go)
- adapters/standard/standard.go
  - [example](./tree/adapters/standard/standard.go/example.go)

- cmd/main/main.go
  - [example](./tree/cmd/main/main.go/example.go)




### Sandbox Changes:
move all hardcoded text  , such as helps,etc, to /assets/ dir
- setup in assets/version.txt file the version of the binary



### Doc Changes
- add documentations n the lib usage section about these new embed mechanic




