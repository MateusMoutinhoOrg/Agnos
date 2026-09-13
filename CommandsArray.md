

### New Tree:

- [sandbox/api/command.go](#sandboxapicommandgo)
- [sandbox/api/sandbox.go](#sandboxapisandboxgo) 
- [sandbox/binds/cli.go](#sandboxbindscligo)
- [sandbox/internal/commands/<item>/new.go](#sandboxinternalcommandsnewgo)
- [sandbox/internal/commands/<item>/entries.yaml](#sandboxinternalcommandsitementriesyaml)
- [sandbox/internal/commands/<item>/handler.go](#sandboxinternalcommandsitemhandlergo)

---
### sandbox/api/command.go

#### action: 
create

#### description: 
Go struct containing all commands proprieds, such as name , paths, flags ,etc, and a handler function of this command.

#### sample:
```go

type CommandArg struct {
    Type        string
    Id          string
    Required    bool
    Description string
    Examples    []string
    Default     string 
}
type CommandFlag struct {
    Type        string
    Id          string
    Required    bool
    Description string
    Examples    []string
    Default     string
    Identifiers []string
}


type Command struct {
	Name        string
    Category    string
    Help        string
    LongDescription string
	Args        []CommandArg
	Flags       []CommandFlag

    GetItem     func(id string) []any
	Handler     func() int 
}
```


---
### sandbox/api/sandbox.go 
#### action: 
    modify
#### modification:
    create the prop commands (a array of Command)


---
### sandbox/binds/cli.go
#### action: 
    modify
#### modification:
    add the constructions of api.commands, by importing and calling each NewCommand() of sandbox/internal/commands/<command>/command.gos

---
### sandbox/internal/commands/<item>/entries.go
#### action: 
    remove 
#### reason:
    its will not be nescessary anymore.


---
### sandbox/internal/commands/<item>/new.go
#### action:
    create
#### description:
    a function that returns a instance of commands constructing only the props, based on entries.yaml.
#### sample:
```go
func NewCommand(sandbox *api.Sandbox) api.Command {
    command := api.NewCommand() // vanila constructions containing the basic methods.
    command.name = "<item>"
    
    command.Args = []api.CommandArg{
        {
            Id: "<item>",
            Required: true,
            Description: "<item>",
            Examples: []string{"<item>"},
            Default: "<item>",
        },
    }

    command.Flags = []api.CommandFlag{
        {
            Id: "<item>",
            Required: true,
            Description: "<item>",
            Examples: []string{"<item>"},
            Default: "<item>",
            Identifiers: []string{"<item>"},
        },
    }
    command.Handler = func() int {
        // the handler recives sandbox and self, to be able to retrive itens using the GetItemMethod.
        return CommandHandler(sandbox,&command)
    }
    return &command

}
```






---
### sandbox/internal/commands/<item>/entries.yaml
#### action:
    nothing, entries.yaml can keep as it is, but can be modified if nescessary
---
### sandbox/internal/commands/<item>/handler.go
#### action:
    modify
#### modificatioN:
    modify the assignature of the functions.to be equal as it colled on. new.go