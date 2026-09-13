

### New Tree:

- [sandbox/api/commands.go](#sandboxapicommandsgo)
- [sandbox/api/sandbox.go](#sandboxapisandboxgo) 
- [sandbox/binds/cli.go](#sandboxbindscli)
- [sandbox/internal/commands/<item>/entries.go](#sandboxinternalcommandsitementriesgo)
- [sandbox/internal/commands/<item>/commands.go](#sandboxinternalcommandsitemcommandsgo)
- [sandbox/internal/commands/<item>/entries.yaml](#sandboxinternalcommandsitementriesyaml)
- [sandbox/internal/commands/<item>/handler.go](#sandboxinternalcommandsitemhandlergo)

### sandbox/api/commands.go


#### description: 
Go struct containing all commands proprieds, such as name , paths, flags ,etc, and a handler function of this command.

#### sample:
```go



type CommandArg struct {
    Name        string
    Type        string
    Required    bool
    Description string
    Examples    []string
    Default     string 
}
type CommandFlag struct {
    Name        string
    Type        string
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


	Handler     func() ()
}
```


### sandbox/api/sandbox.go 

### sandbox/binds/cli.go

### sandbox/internal/commands/<item>/entries.go

### sandbox/internal/commands/<item>/commands.go

### sandbox/internal/commands/<item>/entries.yaml

### sandbox/internal/commands/<item>/handler.go

# 