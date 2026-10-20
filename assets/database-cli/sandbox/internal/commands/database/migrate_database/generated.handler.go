package migrate_database

import (
	"{{.Module}}/sandbox/api"
	opinionatedagnosdatabase "{{.Module}}/sandbox/deps/OpinionatedAgnosDatabase"
	"{{.Module}}/sandbox/internal/commandprops"
	"{{.Module}}/sandbox/internal/databases"
)

// Handle migrates every database the project declares, or the one named, up
// to the version its database.yaml declares. The engine is the
// OpinionatedAgnosDatabase lib's; what each database is reads off its
// generated.migrations.go. It reports each step on stderr and answers with
// the first failure, which names the database and the step.
func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	migrations := []opinionatedagnosdatabase.Migrations{}
	for _, database := range databases.Migrations(sandbox) {
		if input.Database == "" || input.Database == database.Name {
			migrations = append(migrations, database)
		}
	}
	if len(migrations) == 0 {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", "the project declares no database named "+input.Database)
	}

	err := sandbox.Deps.OpinionatedAgnosDatabase.Migrate(opinionatedagnosdatabase.MigrateProps{
		Databases:  sandbox.Deps.DatabaseDeps.Databases,
		IoDeps:     sandbox.Deps.IoDeps,
		Logf:       response.Logf,
		Migrations: migrations,
	})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
