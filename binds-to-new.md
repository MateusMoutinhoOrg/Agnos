

### remove:
- sandbox/binds  


___
### Tree:
- [sandbox/internal/actions/new.go](#sandboxinternalactionsnewgo)
- [sandbox/internal/cli/new.go](#sandboxinternalclinewgo)
- [sandbox/new.go](#sandboxnewgo)
___
### sandbox/internal/actions/new.go
#### sample:
~~~go

func NewActions(sandbox *api.Sandbox) api.Actions {
	actions := api.Actions{}
    actions.Build = func(props api.BuildProps) error {
		return buildAction.Build(sandbox, props)
	}

    return actions
}



~~~



___
### sandbox/internal/cli/new.go
#### sample:
~~~go

func NewCli(sandbox *api.Sandbox) api.Cli {
	cli := api.Cli{}
    cli.Commands = []api.Command{
        new.NewCommand(sandbox),
        deps_init.NewCommand(sandbox),
        deps_purge.NewCommand(sandbox),
        front_init.NewCommand(sandbox),
        front_purge.NewCommand(sandbox),
        server_init.NewCommand(sandbox),
        server_purge.NewCommand(sandbox),
        cli_init.NewCommand(sandbox),
        cli_purge.NewCommand(sandbox),
        build.NewCommand(sandbox),
        compile.NewCommand(sandbox),
        exec_test.NewCommand(sandbox),
        update_test.NewCommand(sandbox),
        verify.NewCommand(sandbox),
        version.NewCommand(sandbox),
        help.NewCommand(sandbox),
    }
    
    return cli
}
~~~
___
### sandbox/new.go
#### sample:
~~~go

func New(deps *deps.Deps) *api.Sandbox {
	self := api.Sandbox{Deps: deps}

	self.Actions = NewActions(&self)
	self.Cli = NewCli(&self)

	return &self
}

~~~