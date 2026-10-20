package database_version

import (
	"{{.Module}}/sandbox/api"
	opinionatedagnosdatabase "{{.Module}}/sandbox/deps/OpinionatedAgnosDatabase"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/databases"
)

// unguarded are the commands a database at another version does not stop:
// the one that migrates it, and the two that read nothing of it.
var unguarded = map[string]bool{"migrate-database": true, "help": true, "version": true}

// Handle refuses a command line when one of the project's databases is not at
// the version its database.yaml declares, naming the command that migrates
// it. It runs after database_dir, so it checks the folder --database named.
// It answers nothing otherwise, so the command after it runs.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	if len(input.Verb) == 0 || unguarded[input.Verb[0]] {
		return nil
	}

	err := sandbox.Deps.OpinionatedAgnosDatabase.CheckVersions(opinionatedagnosdatabase.MigrateProps{
		Databases:  sandbox.Deps.DatabaseDeps.Databases,
		IoDeps:     sandbox.Deps.IoDeps,
		Command:    sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName) + " migrate-database",
		Migrations: databases.Migrations(sandbox),
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	return nil
}
