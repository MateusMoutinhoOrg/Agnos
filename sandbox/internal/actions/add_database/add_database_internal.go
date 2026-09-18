package add_database

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDatabaseInternal writes the one hand-declared file of a new database
// package. It refuses to overwrite an existing declaration (via io.WriteFile).
// The prefix is the key every record is written under; left out, it is the
// database's own name, so two databases of one project never share a keyspace
// by accident.
func AddDatabaseInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string, prefix string) error {
	if err := utils.RequireExtension(sandbox, io, utils.ExtensionSandboxDatabase); err != nil {
		return err
	}

	if err := utils.ValidateDatabaseName(sandbox, name); err != nil {
		return err
	}
	utils.NoteNormalizedCommandName(sandbox, name)

	conf := databaseconf.NewEmpty(sandbox)
	conf.Name = utils.DatabaseIdentifier(sandbox, name)
	conf.Prefix = sandbox.Deps.Stringsdeps.TrimSpace(prefix)
	if conf.Prefix == "" {
		conf.Prefix = conf.Name
	}

	sandbox.Deps.Std.Log("add-database creating %s \n", utils.DatabaseDir(sandbox, name))

	return io.WriteFile(utils.DatabaseConfPath(sandbox, name), []byte(conf.Render()))
}
