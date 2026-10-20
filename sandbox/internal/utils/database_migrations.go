package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// The helpers below are how a database.yaml's version moves. Every edit that
// changes what the stored data holds opens a version of its own: it freezes
// the declaration as it was into migrations/v<N>.yaml, writes what brings the
// data of version N to N+1 into migrations/v<N+1>.ops.yaml, and raises the
// declaration to N+1. The OpinionatedAgnosDatabase lib runs those files, read
// back by the build into each database's generated.migrations.go.

// DatabaseMigrationsDir is the folder of a database package holding one frozen
// schema per past version, the operations of every version past the first,
// and the hand-written step of any version that has one.
const DatabaseMigrationsDir = "migrations"

// The names the migration engine reserves, the OpinionatedAgnosDatabase lib's
// own: the table every database records the version of its data in, and the
// two folders a migration works in inside the --database folder.
const (
	DatabaseVersionTable    = "agnos-version"
	DatabaseMigratingDir    = "agnos-migrating"
	DatabasePreMigrationDir = "agnos-pre-migration"
)

// DatabaseMigrationsPath is the project-relative migrations folder of a
// database.
func DatabaseMigrationsPath(sandbox *api.Sandbox, name string) string {
	return DatabaseDir(sandbox, name) + "/" + DatabaseMigrationsDir
}

// DatabaseSchemaPath is the frozen declaration of one past version.
func DatabaseSchemaPath(sandbox *api.Sandbox, name string, version int) string {
	return DatabaseMigrationsPath(sandbox, name) + "/" + DatabaseVersionFile(sandbox, version, ".yaml")
}

// DatabaseOpsPath is what brings the data to one version from the one before.
func DatabaseOpsPath(sandbox *api.Sandbox, name string, version int) string {
	return DatabaseMigrationsPath(sandbox, name) + "/" + DatabaseVersionFile(sandbox, version, ".ops.yaml")
}

// DatabaseStepPath is the hand-written step of one version.
func DatabaseStepPath(sandbox *api.Sandbox, name string, version int) string {
	return DatabaseMigrationsPath(sandbox, name) + "/" + DatabaseVersionFile(sandbox, version, ".go")
}

// DatabaseVersionFile is the name of one file of a migrations folder:
// v<N> and its suffix.
func DatabaseVersionFile(sandbox *api.Sandbox, version int, suffix string) string {
	return "v" + sandbox.Deps.StringsDeps.FormatInt(int64(version), 10) + suffix
}

// DatabaseStepFunction is the Go name of one version's hand-written step.
func DatabaseStepFunction(sandbox *api.Sandbox, version int) string {
	return "V" + sandbox.Deps.StringsDeps.FormatInt(int64(version), 10)
}

// OpenDatabaseVersion raises a database to the next version: before — the
// declaration as it was rendered ahead of the edit — is frozen as the current
// version, ops become the next one's operations, and conf, already edited,
// takes the next number. The caller saves conf. Both files are created, never
// overwritten: a version once opened is history.
func OpenDatabaseVersion(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, before string,
	conf *databaseconf.DatabaseConf, ops []migrationconf.Op) error {

	current := conf.Version
	if current < 1 {
		current = 1
	}
	next := current + 1

	if err := io.CreateFile(DatabaseSchemaPath(sandbox, name, current), []byte(before)); err != nil {
		return sandbox.Deps.StdDeps.Errorf("cannot freeze version %d of %s: %w", current, DatabaseName(sandbox, name), err)
	}

	ops_conf := migrationconf.NewEmpty(sandbox)
	ops_conf.Ops = ops
	if err := io.CreateFile(DatabaseOpsPath(sandbox, name, next), []byte(ops_conf.Render())); err != nil {
		return sandbox.Deps.StdDeps.Errorf("cannot open version %d of %s: %w", next, DatabaseName(sandbox, name), err)
	}

	conf.Version = next

	how := "migrate it through databases.Migrations (docs/Databases)"
	if has_cli, _ := ExtensionEnabled(sandbox, io, ExtensionCli); has_cli {
		if project_conf, err := LoadProjectConf(sandbox, io); err == nil {
			how = "run `" + sandbox.Deps.StringsDeps.ToLower(project_conf.ProjectName) + " migrate-database`"
		}
	}
	sandbox.Deps.StdDeps.Logf("%s now expects version %d of its data: %s wherever it already holds some \n",
		DatabaseName(sandbox, name), next, how)
	return nil
}

// WriteDatabaseStep writes the hand-written step of one version, once: a step
// already on disk is the project's and is left as it is.
func WriteDatabaseStep(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, version int) error {
	dest := DatabaseStepPath(sandbox, name, version)
	if io.IsFile(dest) {
		return nil
	}

	module_conf, err := LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("writing the hand-written step of version %d in %s \n", version, dest)

	return RenderTemplateToDest(sandbox, io, "templates/database_migration_step.go", map[string]any{
		"Module":        module_conf.Module,
		"GeneratorName": sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName),
		"DatabaseName":  DatabaseName(sandbox, name),
		"Version":       version,
		"Previous":      version - 1,
		"Function":      DatabaseStepFunction(sandbox, version),
	}, dest)
}

// DatabaseFieldEditOps is what one edit of a field changes about the data
// already stored, in the order it runs: the values move to the new name, then
// read as the new type, then a field that became required is filled. An edit
// that changes nothing stored — a field made optional, a link pointed at
// another table, a type every stored value already reads as — has none.
func DatabaseFieldEditOps(table string, parent string, current databaseconf.Field, edited databaseconf.Field, fill string) []migrationconf.Op {
	ops := []migrationconf.Op{}
	name := current.Name

	if edited.Name != current.Name {
		ops = append(ops, migrationconf.Op{Kind: migrationconf.OpRenameField, Table: table, Parent: parent,
			Field: current.Name, Type: current.Type, To: edited.Name})
		name = edited.Name
	}

	reshaped := (current.Type == databaseconf.FieldObject) != (edited.Type == databaseconf.FieldObject)
	if reshaped {
		ops = append(ops, migrationconf.Op{Kind: migrationconf.OpRemoveField, Table: table, Parent: parent,
			Field: name, Type: current.Type})
	} else if edited.Type != current.Type && DatabaseRetypeChangesData(current.Type, edited.Type) {
		ops = append(ops, migrationconf.Op{Kind: migrationconf.OpRetypeField, Table: table, Parent: parent,
			Field: name, Type: current.Type, To: edited.Type})
	}

	if edited.Required && (!current.Required || reshaped) {
		ops = append(ops, DatabaseRequireOp(table, parent, edited, fill))
	}

	return ops
}

// DatabaseRequireOp is the operation of a field that became required.
func DatabaseRequireOp(table string, parent string, field databaseconf.Field, fill string) migrationconf.Op {
	return migrationconf.Op{Kind: migrationconf.OpRequireField, Table: table, Parent: parent,
		Field: field.Name, Type: field.Type, Default: fill}
}

// DatabaseRetypeChangesData reports a type change some stored value may not
// read as: one into a key, which has to be indexed and held unique, or one
// into a number from text. A value written as text reads back as text
// whatever it held, and a whole number as a number, so those need no version.
func DatabaseRetypeChangesData(from string, to string) bool {
	switch to {
	case databaseconf.FieldString, databaseconf.FieldBytes:
		return false
	case databaseconf.FieldInteger, databaseconf.FieldLink:
		return from != databaseconf.FieldInteger && from != databaseconf.FieldLink
	case databaseconf.FieldNumber:
		return from != databaseconf.FieldInteger && from != databaseconf.FieldLink && from != databaseconf.FieldNumber
	}
	return true
}

// CheckDatabaseDefault holds a --default to the type of the field it fills,
// so a migration never stops on a value the command line could have refused.
// A key takes none: every record it filled would hold the same unique value.
func CheckDatabaseDefault(sandbox *api.Sandbox, field databaseconf.Field, fill string) error {
	switch field.Type {
	case databaseconf.FieldKey:
		return sandbox.Deps.StdDeps.Errorf("a key takes no --default: every record it filled would hold the same unique value, so fill it by hand (--fill-by-hand)")
	case databaseconf.FieldInteger, databaseconf.FieldLink:
		if _, err := sandbox.Deps.StringsDeps.ParseInt(fill, 10, 64); err != nil {
			return sandbox.Deps.StdDeps.Errorf("--default %q is not a whole number, which a %s field holds", fill, field.Type)
		}
	case databaseconf.FieldNumber:
		if _, err := sandbox.Deps.StringsDeps.ParseFloat(fill, 64); err != nil {
			return sandbox.Deps.StdDeps.Errorf("--default %q is not a number, which a number field holds", fill)
		}
	}
	return nil
}

// ValidateDatabaseKeyPrefix holds a key-prefix to what a migration can copy
// and move: a folder inside the --database folder, its segments the alphabet
// the store keeps as it is on disk, and none of the folders a migration works
// in.
func ValidateDatabaseKeyPrefix(sandbox *api.Sandbox, prefix string) error {
	segments := sandbox.Deps.StringsDeps.Split(sandbox.Deps.StringsDeps.TrimSuffix(prefix, "/"), "/")
	for _, segment := range segments {
		if segment == "" {
			return sandbox.Deps.StdDeps.Errorf("the key-prefix %q holds an empty segment", prefix)
		}
		for _, letter := range segment {
			valid := (letter >= 'a' && letter <= 'z') || (letter >= '0' && letter <= '9') || letter == '-' || letter == '_'
			if !valid {
				return sandbox.Deps.StdDeps.Errorf(
					"the key-prefix %q holds %q: its segments are lower-case letters, digits, - and _, the folder names a migration copies and moves",
					prefix, string(letter))
			}
		}
	}
	if segments[0] == DatabaseMigratingDir || segments[0] == DatabasePreMigrationDir {
		return sandbox.Deps.StdDeps.Errorf("the key-prefix %q starts with %s, a folder migrate-database works in", prefix, segments[0])
	}
	return nil
}
