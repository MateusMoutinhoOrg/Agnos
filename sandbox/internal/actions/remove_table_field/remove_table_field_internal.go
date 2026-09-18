package remove_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveTableFieldInternal parses the target database's specs.yaml, drops the
// named field from the table — or from the nested collection --parent names —
// and writes the file back.
func RemoveTableFieldInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.DatabaseFieldProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	name := utils.DatabaseName(sandbox, props.Name)
	index := utils.FindDatabaseField(sandbox, fields, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("%s declares no field %q", where, name)
	}

	sandbox.Deps.Std.Log("remove-table-field removing %s.%s from %s \n", where, name, utils.DatabaseConfPath(sandbox, props.Database))

	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, utils.RemoveDatabaseField(fields, index)); err != nil {
		return err
	}
	return utils.SaveDatabaseConf(sandbox, io, props.Database, conf)
}
