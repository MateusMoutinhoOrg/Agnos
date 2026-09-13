

### New Tree:

- [sandbox/api/commands.go](#sandboxapicommandsgo)
- [sandbox/api/sandbox.go](#sandboxapisandboxgo) 
- [sandbox/binds/cli.go](#sandboxbindscli)
- [sandbox/internal/commands/<item>/new.go](#sandboxinternalcommandsnew.go)
- [sandbox/internal/commands/<item>/entries.yaml](#sandboxinternalcommandsitementriesyaml)
- [sandbox/internal/commands/<item>/handler.go](#sandboxinternalcommandsitemhandlergo)

---
### sandbox/api/commands.go

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
### sandbox/internal/commands/<item>/new-command.go
#### action:
create
#### description:
a function that returns a instance of commands constructing only the props, based on entries.yaml.
#### sample:
```go
func NewCommand() api.Command {
    command := NewBasicCommand()
    command.name = "<item>"
    

    return api.Command{
        Name: "<item>",
        Category: "<item>",
        Help: "<item>",
        LongDescription: "<item>",
        Args: []api.CommandArg{
            {
                Id: "<item>",
                Required: true,
                Description: "<item>",
                Examples: []string{"<item>"},
                Default: "<item>",
            },
        },
        Flags: []api.CommandFlag{
            {
                Id: "<item>",
                Required: true,
                Description: "<item>",
                Examples: []string{"<item>"},
                Default: "<item>",
                Identifiers: []string{"<item>"},
            },
        },
      
        Handler: func() int {
            return 0
        },
    }
}







---
### sandbox/internal/commands/<item>/entries.yaml


---
### sandbox/internal/commands/<item>/handler.go

# 