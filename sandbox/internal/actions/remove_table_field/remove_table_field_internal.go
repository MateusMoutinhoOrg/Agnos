package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/migrationconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveTableFieldInternal parses the target database's database.yaml, drops the
// named field from the table — or from the nested collection --parent names —
// and writes the file back. It opens a version whose one operation deletes the
// values every record already stored holds for it, so a field declared again
// under that name later never reads an old one.
func RemoveTableFieldInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.RemoveTableFieldProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}
	before := conf.Render()

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	name := utils.DatabaseName(sandbox, props.Name)
	index := utils.FindDatabaseField(sandbox, fields, name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("%s declares no field %q", where, name)
	}

	sandbox.Deps.StdDeps.Logf("remove-table-field removing %s.%s from %s \n", where, name, utils.DatabaseConfPath(sandbox, props.Database))

	op := migrationconf.Op{Kind: migrationconf.OpRemoveField, Table: utils.DatabaseName(sandbox, props.Table),
		Parent: utils.DatabaseName(sandbox, props.Parent), Field: name, Type: fields[index].Type}

	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, utils.RemoveAt(fields, index)); err != nil {
		return err
	}
	if err := utils.OpenDatabaseVersion(sandbox, io, props.Database, before, conf, []migrationconf.Op{op}); err != nil {
		return err
	}
	return utils.SaveDatabaseConf(sandbox, io, props.Database, conf)
}
