package add_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddTableInternal parses the target database's database.yaml, appends the new
// table and writes the file back. A table is born with no fields: what it
// holds is declared one `add-table-field` at a time.
func AddTableInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, database string, table string) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return err
	}

	if err := utils.ValidateDatabaseMember(sandbox, "table", table); err != nil {
		return err
	}
	utils.NoteNormalizedName(sandbox, "table", table)

	name := utils.DatabaseName(sandbox, table)
	if utils.FindDatabaseTable(sandbox, conf.Tables, name) >= 0 {
		return sandbox.Deps.StdDeps.Errorf("database %q already declares the table %q",
			utils.DatabaseName(sandbox, database), name)
	}

	sandbox.Deps.StdDeps.Logf("add-table adding %s to %s \n", name, utils.DatabaseConfPath(sandbox, database))

	conf.Tables = append(conf.Tables, databaseconf.Table{Name: name, Fields: []databaseconf.Field{}})
	return utils.SaveDatabaseConf(sandbox, io, database, conf)
}
