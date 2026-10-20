package add_database_migration

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDatabaseMigrationInternal parses the target database's database.yaml,
// opens its next version with no operation of its own — the schema does not
// change, the data does — writes that version's hand-written step once and
// writes the declaration back. It is the one version no editing command
// opens: a seed, a field computed from others, a value split in two.
func AddDatabaseMigrationInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, database string) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("add-database-migration opening version %d of %s \n", conf.Version+1, utils.DatabaseConfPath(sandbox, database))

	if err := utils.OpenDatabaseVersion(sandbox, io, database, conf.Render(), conf, []migrationconf.Op{}); err != nil {
		return err
	}
	if err := utils.WriteDatabaseStep(sandbox, io, database, conf.Version); err != nil {
		return err
	}
	return utils.SaveDatabaseConf(sandbox, io, database, conf)
}
