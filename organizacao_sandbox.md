Troque a logica de funcionamento da lib do sandbox, para ao ivez de estender elementos, ela ter subpropriedades


## Nova Estrutura: 
- contracts/lib/cli/cli.go
- contracts/lib/config/config.go
- contracts/lib/core/core.go
- contracts/sandbox/sandbox.go

### Estrutura do sandbox: 
~~~go


type SandBox struct {
	cli CliApi 
	config Config
	core CoreApi
	Deps deps.Deps
}

~~~