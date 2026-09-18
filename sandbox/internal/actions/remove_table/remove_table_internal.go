package remove_table

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveTableInternal parses the target database's specs.yaml, drops the named
// table and writes the file back. A table another table links to is refused:
// the link would name a collection the declaration no longer has, which is the
// one shape the generator cannot render.
func RemoveTableInternal(sandbox *api.Sandbox, io *smartio.SmartIO, database string, table string) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, database)
	if err != nil {
		return err
	}

	name := utils.DatabaseName(sandbox, table)
	index := utils.FindDatabaseTable(sandbox, conf.Tables, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("database %q declares no table %q",
			utils.DatabaseIdentifier(sandbox, database), name)
	}

	for _, other := range conf.Tables {
		if other.Name == name {
			continue
		}
		for _, field := range other.Fields {
			if field.Type == databaseconf.FieldLink && field.Target == name {
				return sandbox.Deps.Std.Errorf(
					"table %q is the target of the link %s.%s: drop that field first, or point it elsewhere",
					name, other.Name, field.Name)
			}
		}
	}

	sandbox.Deps.Std.Log("remove-table removing %s from %s \n", name, utils.DatabaseConfPath(sandbox, database))

	conf.Tables = utils.RemoveDatabaseTable(conf.Tables, index)
	return utils.SaveDatabaseConf(sandbox, io, database, conf)
}
