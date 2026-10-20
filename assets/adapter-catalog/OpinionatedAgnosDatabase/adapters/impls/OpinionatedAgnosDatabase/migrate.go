package opinionatedagnosdatabase

import (
	"fmt"
	"strconv"
	"strings"

	opinionatedagnosdatabase "{{.Module}}/sandbox/deps/OpinionatedAgnosDatabase"
	"{{.Module}}/sandbox/deps/databasedeps"
)

// The migration engine. A database's folder is copied into MigratingDir, every
// step from the version its data records runs on the copy, the copy is
// repaired and checked against the last schema, and only then is the live
// folder moved into PreMigrationDir and the copy moved in its place. The
// store has no transaction across keys, so the copy is what makes a step that
// stops half way harmless: the live folder is only ever read until the swap,
// and the swap is two renames a crash between is finished by the next run.

// versionField is the one field of VersionTable.
const versionField = "version"

// doneSuffix marks a copy that passed every step and check: beside it, it is
// what tells a crash during the swap, which the next run finishes, from one
// during a step, whose copy is thrown away. No key of the store holds a dot,
// so no record can collide with it.
const doneSuffix = ".done"

// folders is where one database lives on disk and where a migration of it
// works, as host paths and as the Path the store opens each of them under.
type folders struct {
	live     string
	work     string
	done     string
	keep     string
	livePath string
	workPath string
}

// foldersOf is the folders of one database: a key-prefix is a folder inside
// the --database folder, and a migration's copy is the same path inside
// MigratingDir.
func foldersOf(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations) folders {
	io := props.IoDeps
	work := io.Join(database.Dir, opinionatedagnosdatabase.MigratingDir, database.KeyPrefix)
	return folders{
		live:     io.Join(database.Dir, database.KeyPrefix),
		work:     work,
		done:     work + doneSuffix,
		keep:     io.Join(database.Dir, opinionatedagnosdatabase.PreMigrationDir, database.KeyPrefix),
		livePath: database.Dir + "/" + database.KeyPrefix,
		workPath: database.Dir + "/" + opinionatedagnosdatabase.MigratingDir + "/" + database.KeyPrefix,
	}
}

// migrate is Contract.Migrate.
func migrate(props opinionatedagnosdatabase.MigrateProps) error {
	for _, database := range props.Migrations {
		if err := migrateOne(props, database); err != nil {
			return fmt.Errorf("%s: %w", database.Name, err)
		}
	}
	return nil
}

// migrateOne brings one database up to its last version.
func migrateOne(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations) error {
	if len(database.Versions) == 0 {
		return nil
	}
	last := database.Versions[len(database.Versions)-1].Version
	at := foldersOf(props, database)

	if err := recoverMigration(props, database, at); err != nil {
		return err
	}

	stored, known, err := storedVersion(props.Databases, at.livePath)
	if err != nil {
		return err
	}
	if !known {
		empty, err := isEmpty(props.Databases, at.livePath, database)
		if err != nil {
			return err
		}
		if empty {
			logf(props, "%s holds no data: stamped at version %d\n", database.Name, last)
			return writeVersion(props.Databases, at.livePath, last)
		}
		stored = 1
	}

	if stored == last {
		logf(props, "%s is at version %d: nothing to migrate\n", database.Name, last)
		return nil
	}
	if stored > last {
		return fmt.Errorf("its data is at version %d, newer than the version %d this build knows: it was migrated by a newer build, which is the one to run", stored, last)
	}

	if err := props.IoDeps.Copy(at.live, at.work); err != nil {
		props.IoDeps.RemoveDir(at.work)
		return fmt.Errorf("copying %s: %w", at.live, err)
	}
	if err := runSteps(props, database, at, stored); err != nil {
		props.IoDeps.RemoveDir(at.work)
		removeIfEmpty(props, props.IoDeps.Join(database.Dir, opinionatedagnosdatabase.MigratingDir))
		return fmt.Errorf("%w (the data in %s was not touched)", err, at.live)
	}
	if err := props.IoDeps.WriteFile(at.done, []byte(strconv.Itoa(last))); err != nil {
		props.IoDeps.RemoveDir(at.work)
		return fmt.Errorf("marking the copy done: %w (the data in %s was not touched)", err, at.live)
	}
	return swap(props, database, at)
}

// runSteps runs every version past stored on the copy, recording each one
// reached in the copy itself, then repairs and checks the copy against the
// last schema.
func runSteps(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations, at folders, stored int) error {
	for index, version := range database.Versions {
		if index == 0 || version.Version <= stored {
			continue
		}
		logf(props, "%s: migrating to version %d\n", database.Name, version.Version)

		for _, op := range version.Ops {
			if err := applyOp(props, at, op); err != nil {
				return fmt.Errorf("version %d, %s %s: %w", version.Version, op.Kind, opPlace(op), err)
			}
		}

		if version.Custom != nil {
			from, err := open(props.Databases, at.workPath, database.Versions[index-1].Schemas)
			if err != nil {
				return err
			}
			to, err := open(props.Databases, at.workPath, version.Schemas)
			if err != nil {
				return err
			}
			if err := version.Custom(from, to); err != nil {
				return fmt.Errorf("version %d: %w", version.Version, err)
			}
		}

		if err := writeVersion(props.Databases, at.workPath, version.Version); err != nil {
			return err
		}
	}

	last := database.Versions[len(database.Versions)-1]
	return checkData(props.Databases, at.workPath, last)
}

// swap moves the live folder into PreMigrationDir and the finished copy in its
// place. Both moves are renames, and a run that finds a finished copy calls
// it again, so a crash anywhere in it is finished rather than undone.
func swap(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations, at folders) error {
	io := props.IoDeps
	if io.Exists(at.live) {
		previous, known, err := storedVersion(props.Databases, at.livePath)
		if err != nil {
			return err
		}
		if !known {
			previous = 1
		}
		keep := freeName(props, at.keep+"-v"+strconv.Itoa(previous))
		if err := io.Rename(at.live, keep); err != nil {
			return fmt.Errorf("moving %s aside: %w (run the migration again to finish it)", at.live, err)
		}
		logf(props, "%s: the data of version %d is kept in %s\n", database.Name, previous, keep)
	}
	if err := io.Rename(at.work, at.live); err != nil {
		return fmt.Errorf("moving the migrated copy into %s: %w (run the migration again to finish it)", at.live, err)
	}
	io.RemoveDir(at.done)
	removeIfEmpty(props, io.Join(database.Dir, opinionatedagnosdatabase.MigratingDir))
	logf(props, "%s: migrated to version %d\n", database.Name, database.Versions[len(database.Versions)-1].Version)
	return nil
}

// recoverMigration settles what an earlier run left in MigratingDir: a copy
// marked done is swapped in, and any other copy — a run that stopped during a
// step — is thrown away, since the live folder was never touched.
func recoverMigration(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations, at folders) error {
	io := props.IoDeps
	if !io.Exists(at.work) {
		io.RemoveDir(at.done)
		return nil
	}
	if !io.IsFile(at.done) {
		logf(props, "%s: throwing away a migration that stopped before it finished\n", database.Name)
		io.RemoveDir(at.work)
		return nil
	}
	logf(props, "%s: finishing a migration that stopped while swapping its folders\n", database.Name)
	return swap(props, database, at)
}

// checkVersions is Contract.CheckVersions.
func checkVersions(props opinionatedagnosdatabase.MigrateProps) error {
	var refused []string
	for _, database := range props.Migrations {
		if reason := checkOne(props, database); reason != "" {
			refused = append(refused, database.Name+": "+reason)
		}
	}
	if len(refused) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(refused, "\n"))
}

// checkOne is why one database cannot be used by this build, "" when it can.
func checkOne(props opinionatedagnosdatabase.MigrateProps, database opinionatedagnosdatabase.Migrations) string {
	if len(database.Versions) == 0 {
		return ""
	}
	last := database.Versions[len(database.Versions)-1].Version
	at := foldersOf(props, database)

	if props.IoDeps.Exists(at.work) {
		return fmt.Sprintf("a migration of it stopped half way: run `%s` to settle it", props.Command)
	}

	stored, known, err := storedVersion(props.Databases, at.livePath)
	if err != nil {
		return err.Error()
	}
	if !known {
		empty, err := isEmpty(props.Databases, at.livePath, database)
		if err != nil {
			return err.Error()
		}
		if empty {
			if err := writeVersion(props.Databases, at.livePath, last); err != nil {
				return err.Error()
			}
			return ""
		}
		stored = 1
	}

	if stored < last {
		return fmt.Sprintf("its data is at version %d and this build expects version %d: run `%s`", stored, last, props.Command)
	}
	if stored > last {
		return fmt.Sprintf("its data is at version %d, newer than the version %d this build knows: it was migrated by a newer build, which is the one to run", stored, last)
	}
	return ""
}

// storedVersion is the version a database's data records, and false when it
// records none.
func storedVersion(databases databasedeps.Databases, path string) (int, bool, error) {
	record, found, err := versionRecord(databases, path)
	if err != nil || !found {
		return 0, false, err
	}
	raw, failure := record.Get(versionField)
	if failure != nil {
		return 0, false, fail(failure)
	}
	version, ok := raw.(int64)
	if !ok {
		return 0, false, mistyped(versionField, "a whole number")
	}
	return int(version), true, nil
}

// writeVersion records the version a database's data is at.
func writeVersion(databases databasedeps.Databases, path string, version int) error {
	record, found, err := versionRecord(databases, path)
	if err != nil {
		return err
	}
	if found {
		return fail(record.Update(versionField, int64(version)))
	}
	database, err := open(databases, path, versionSchemas())
	if err != nil {
		return err
	}
	versions, err := collection(database, opinionatedagnosdatabase.VersionTable)
	if err != nil {
		return err
	}
	_, failure := versions.Insert(map[string]any{versionField: int64(version)})
	return fail(failure)
}

// versionRecord is the one record of VersionTable.
func versionRecord(databases databasedeps.Databases, path string) (databasedeps.Record, bool, error) {
	database, err := open(databases, path, versionSchemas())
	if err != nil {
		return databasedeps.Record{}, false, err
	}
	versions, err := collection(database, opinionatedagnosdatabase.VersionTable)
	if err != nil {
		return databasedeps.Record{}, false, err
	}
	records, failure := versions.List(1, 1)
	if failure != nil {
		return databasedeps.Record{}, false, fail(failure)
	}
	if len(records) == 0 {
		return databasedeps.Record{}, false, nil
	}
	return records[0], true, nil
}

// versionSchemas is the schema VersionTable is opened with.
func versionSchemas() []databasedeps.Schema {
	version := databasedeps.Field{Name: versionField, Type: databasedeps.Int, Required: true}
	table := databasedeps.Schema{Name: opinionatedagnosdatabase.VersionTable, Fields: []databasedeps.Field{version}}
	return []databasedeps.Schema{table}
}

// isEmpty reports a database holding no record in any table any of its
// versions declared.
func isEmpty(databases databasedeps.Databases, path string, database opinionatedagnosdatabase.Migrations) (bool, error) {
	seen := map[string]bool{}
	for _, version := range database.Versions {
		for _, schema := range version.Schemas {
			if seen[schema.Name] {
				continue
			}
			seen[schema.Name] = true

			bare, err := open(databases, path, []databasedeps.Schema{databasedeps.Schema{Name: schema.Name}})
			if err != nil {
				return false, err
			}
			table, err := collection(bare, schema.Name)
			if err != nil {
				return false, err
			}
			records, failure := table.List(1, 1)
			if failure != nil {
				return false, fail(failure)
			}
			if len(records) > 0 {
				return false, nil
			}
		}
	}
	return true, nil
}

// open builds a database over path with the given schemas.
func open(databases databasedeps.Databases, path string, schemas []databasedeps.Schema) (databasedeps.Database, error) {
	database, failure := databases.New(databasedeps.Props{Path: path, Schemas: schemas})
	if failure != nil {
		return databasedeps.Database{}, fail(failure)
	}
	return database, nil
}

// freeName is base, or base-2, base-3… when an earlier migration already
// kept a folder under it: a kept folder is never written over.
func freeName(props opinionatedagnosdatabase.MigrateProps, base string) string {
	if !props.IoDeps.Exists(base) {
		return base
	}
	for index := 2; ; index++ {
		candidate := base + "-" + strconv.Itoa(index)
		if !props.IoDeps.Exists(candidate) {
			return candidate
		}
	}
}

// removeIfEmpty removes a folder left holding nothing.
func removeIfEmpty(props opinionatedagnosdatabase.MigrateProps, path string) {
	if props.IoDeps.IsDir(path) && len(props.IoDeps.ListAll(path)) == 0 {
		props.IoDeps.RemoveDir(path)
	}
}

// logf reports one step, when the caller gave a channel to report on.
func logf(props opinionatedagnosdatabase.MigrateProps, format string, args ...any) {
	if props.Logf != nil {
		props.Logf(format, args...)
	}
}
