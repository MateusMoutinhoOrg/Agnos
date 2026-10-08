package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveTableInternal parses the target database's database.yaml, drops the named
// table and writes the file back. A table another table links to is refused:
// the link would name a collection the declaration no longer has, which is the
// one shape the generator cannot render.
func RemoveTableInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, database string, table string) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return err
	}

	name := utils.DatabaseName(sandbox, table)
	index := utils.FindDatabaseTable(sandbox, conf.Tables, name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("database %q declares no table %q",
			utils.DatabaseName(sandbox, database), name)
	}

	for _, other := range conf.Tables {
		if other.Name == name {
			continue
		}
		for _, field := range other.Fields {
			if field.Type == databaseconf.FieldLink && field.Target == name {
				return sandbox.Deps.StdDeps.Errorf(
					"table %q is the target of the link %s.%s: drop that field first, or point it elsewhere",
					name, other.Name, field.Name)
			}
		}
	}

	sandbox.Deps.StdDeps.Logf("remove-table removing %s from %s \n", name, utils.DatabaseConfPath(sandbox, database))

	conf.Tables = utils.RemoveAt(conf.Tables, index)
	return utils.SaveDatabaseConf(sandbox, io, database, conf)
}
