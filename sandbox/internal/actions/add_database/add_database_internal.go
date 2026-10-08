package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDatabaseInternal writes the one hand-declared file of a new database
// package. It refuses to overwrite an existing declaration (via io.WriteFile).
// The prefix is the key every record is written under; left out, it is the
// database's own name, so two databases of one project never share a keyspace
// by accident.
func AddDatabaseInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, name string, prefix string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionDatabase); err != nil {
		return err
	}

	if err := utils.ValidateDatabaseName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedName(sandbox, "database", name)

	conf := databaseconf.NewEmpty(sandbox)
	conf.Name = utils.DatabaseName(sandbox, name)
	conf.KeyPrefix = sandbox.Deps.StringsDeps.TrimSpace(prefix)
	if conf.KeyPrefix == "" {
		conf.KeyPrefix = conf.Name
	}

	sandbox.Deps.StdDeps.Logf("add-database creating %s \n", utils.DatabaseDir(sandbox, name))

	return io.CreateFile(utils.DatabaseConfPath(sandbox, name), []byte(conf.Render()))
}
