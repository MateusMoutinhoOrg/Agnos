

### New Tree:

- [sandbox/api/commands.go](#sandboxapicommandsgo)
- [sandbox/api/sandbox.go](#sandboxapisandboxgo) 
- [sandbox/binds/cli.go](#sandboxbindscli)
- [sandbox/internal/commands/<item>/commands.go](#sandboxinternalcommandsitemcommandsgo)
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
    Id          string
    Required    bool
    Description string
    Examples    []string
    Default     string 
}
type CommandFlag struct {
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
add the constructions of api.commands


---
### sandbox/internal/commands/<item>/entries.go
#### action: 
remove 
#### reason:
its will not be nescessary anymore.


---
### sandbox/internal/commands/<item>/commands.go
#### action:
create



---
### sandbox/internal/commands/<item>/entries.yaml


---
### sandbox/internal/commands/<item>/handler.go

# 