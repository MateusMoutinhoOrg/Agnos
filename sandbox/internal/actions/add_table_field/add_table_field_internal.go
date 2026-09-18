package add_table_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddTableFieldInternal parses the target database's specs.yaml, appends the
// new field to the table — or to the nested collection --parent names — and
// writes the file back. Where the field lands is what decides which methods it
// generates: a `key` brings a Find, a `link` a Get, a `database` the Add/List
// pair, and every plain field an Update and a place in the filtrage.
func AddTableFieldInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.DatabaseFieldProps) error {
	conf, err := utils.LoadDatabaseConf(sandbox, io, props.Database)
	if err != nil {
		return err
	}

	fields, where, err := utils.DatabaseFieldsAt(sandbox, conf, props.Table, props.Parent)
	if err != nil {
		return err
	}

	field, err := utils.NewDatabaseField(sandbox, conf, props)
	if err != nil {
		return err
	}

	if utils.FindDatabaseField(sandbox, fields, field.Name) >= 0 {
		return sandbox.Deps.Std.Errorf("%s already declares a field named %q", where, field.Name)
	}

	sandbox.Deps.Std.Log("add-table-field adding %s.%s to %s \n", where, field.Name, utils.DatabaseConfPath(sandbox, props.Database))

	if err := utils.SetDatabaseFieldsAt(sandbox, conf, props.Table, props.Parent, append(fields, field)); err != nil {
		return err
	}
	return utils.SaveDatabaseConf(sandbox, io, props.Database, conf)
}
